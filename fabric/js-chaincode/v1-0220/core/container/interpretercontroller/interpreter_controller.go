/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package interpretercontroller

import (
	"fmt"
	"io"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	container "github.com/hyperledger/fabric/core/container/api"
	"github.com/hyperledger/fabric/core/container/ccintf"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/op/go-logging"

	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
	"path/filepath"
	"github.com/hyperledger/fabric/core/chaincode/platforms/javascript/jscc"
)

type interpreterContainer struct {
	chaincode shim.Chaincode
	running   bool
	args      []string
	env       []string
	stopChan  chan struct{}
}

var (
	interpreterLogger = logging.MustGetLogger("interpretercontroller")
	instRegistry      = make(map[string]*interpreterContainer)
)

//InterpreterVM is a vm. It is identified by chaincode type
type InterpreterVM struct {
}

func (vm *InterpreterVM) getInstance(ccid ccintf.CCID, instName string, args []string, env []string) (*interpreterContainer, error) {
	if ccid.ChaincodeSpec.Type == pb.ChaincodeSpec_JAVASCRIPT {
		ipc := instRegistry[instName]
		if ipc != nil {
			interpreterLogger.Warningf("chaincode instance exists for %s", instName)
			return ipc, nil
		}
		var chaincode shim.Chaincode
		chaincode, err := jscc.NewJavascriptCC(ccid.ChaincodeSpec.ChaincodeId)
		if err != nil {
			return nil, err
		}
		ipc = &interpreterContainer{args: args, env: env, chaincode: chaincode, stopChan: make(chan struct{})}
		instRegistry[instName] = ipc
		interpreterLogger.Debugf("chaincode instance created for %s", instName)
		return ipc, nil
	}
	return nil, fmt.Errorf("unknown script type %d", ccid.ChaincodeSpec.Type)
}

func (vm *InterpreterVM) Deploy(ctxt context.Context, ccid ccintf.CCID, args []string, env []string, reader io.Reader) error {
	//FIXME: deprecated interface??
	return nil
}

func (interpreter *interpreterContainer) launch(ctxt context.Context, id string, args []string, env []string, ccSupport ccintf.CCSupport) error {
	peerRcvCCSend := make(chan *pb.ChaincodeMessage)
	ccRcvPeerSend := make(chan *pb.ChaincodeMessage)
	var err error
	ccchan := make(chan struct{}, 1)
	ccsupportchan := make(chan struct{}, 1)
	go func() {
		defer close(ccchan)
		interpreterLogger.Debugf("chaincode started for %s", id)
		if args == nil {
			args = interpreter.args
		}
		if env == nil {
			env = interpreter.env
		}
		err := shim.StartInProc(env, args, interpreter.chaincode, ccRcvPeerSend, peerRcvCCSend)
		if err != nil {
			err = fmt.Errorf("chaincode-support ended with err: %s", err)
			interpreterLogger.Errorf("%s", err)
		}
		interpreterLogger.Debugf("chaincode ended with for  %s with err: %s", id, err)
	}()

	go func() {
		defer close(ccsupportchan)
		inprocStream := newInProcStream(peerRcvCCSend, ccRcvPeerSend)
		interpreterLogger.Debugf("chaincode-support started for  %s", id)
		err := ccSupport.HandleChaincodeStream(ctxt, inprocStream)
		if err != nil {
			err = fmt.Errorf("chaincode ended with err: %s", err)
			interpreterLogger.Errorf("%s", err)
		}
		interpreterLogger.Debugf("chaincode-support ended with for  %s with err: %s", id, err)
	}()

	select {
	case <-ccchan:
		close(peerRcvCCSend)
		interpreterLogger.Debugf("chaincode %s quit", id)
	case <-ccsupportchan:
		close(ccRcvPeerSend)
		interpreterLogger.Debugf("chaincode support %s quit", id)
	case <-interpreter.stopChan:
		close(ccRcvPeerSend)
		close(peerRcvCCSend)
		interpreterLogger.Debugf("chaincode %s stopped", id)
	}

	return err
}

func (vm *InterpreterVM) Start(ctxt context.Context, ccid ccintf.CCID, args []string, env []string, builder container.BuildSpecFactory) error {
	chaincodeID := ccid.ChaincodeSpec.ChaincodeId
	path := filepath.Join(chaincodeID.Path, chaincodeID.Name + "." + chaincodeID.Version)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		reader, err := builder()
		if err != nil {
			return err
		}
		ioutil.ReadAll(reader)
	}

	instName, _ := vm.GetVMName(ccid)
	interpreter, err := vm.getInstance(ccid, instName, args, env)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("start - could not create interpreter %s", err.Error()))
	}
	if interpreter.running {
		return fmt.Errorf(fmt.Sprintf("chaincode running %s", instName))
	}

	ccSupport, ok := ctxt.Value(ccintf.GetCCHandlerKey()).(ccintf.CCSupport)
	if !ok || ccSupport == nil {
		return errors.New("start - in-process communication generator not supplied")
	}

	interpreterLogger.Debugf("Start container %s", instName)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				interpreterLogger.Criticalf("caught panic from chaincode  %s", instName)
			}
		}()
		interpreter.launch(ctxt, instName, args, env, ccSupport)
	}()

	return nil
}

func (vm *InterpreterVM) Stop(ctxt context.Context, ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {

	instName, _ := vm.GetVMName(ccid)

	ipc := instRegistry[instName]

	if ipc == nil {
		return fmt.Errorf("%s not found", instName)
	}

	if !ipc.running {
		return fmt.Errorf("%s not running", instName)
	}

	ipc.stopChan <- struct{}{}

	delete(instRegistry, instName)
	return nil
}

//Destroy destroys an image
func (vm *InterpreterVM) Destroy(ctxt context.Context, ccid ccintf.CCID, force bool, noprune bool) error {
	//TODO remove script files
	return nil
}

//GetVMName ignores the peer and network name as it just needs to be unique in process
func (vm *InterpreterVM) GetVMName(ccid ccintf.CCID) (string, error) {
	return ccid.GetName(), nil
}
