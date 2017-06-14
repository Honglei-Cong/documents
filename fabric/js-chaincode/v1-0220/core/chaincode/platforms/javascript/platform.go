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

package javascript

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	cutil "github.com/hyperledger/fabric/core/container/util"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func decodeUrl(spec *pb.ChaincodeSpec) (string, error) {
	urlLocation := spec.ChaincodeId.Path
	if urlLocation == "" {
		return "", errors.New("ChaincodeSpec's path/URL cannot be empty")
	}

	if strings.LastIndex(urlLocation, "/") == len(urlLocation)-1 {
		urlLocation = urlLocation[:len(urlLocation)-1]
	}

	return urlLocation, nil
}

//tw is expected to have the chaincode in it from GenerateHashcode. This method
//will just package rest of the bytes
func writeChaincodePackage(spec *pb.ChaincodeSpec, tw *tar.Writer) error {

	urlLocation, err := decodeUrl(spec)
	if err != nil {
		return fmt.Errorf("could not decode url: %s", err)
	}

	err = cutil.WriteFolderToTarPackage(tw, urlLocation, "", map[string]bool{".js": true}, nil)
	if err != nil {
		return fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}
	return nil
}

func getJsChaincodeLocalPath(ccid *pb.ChaincodeID) string {
	return filepath.Join(ccid.Path, ccid.Name+"."+ccid.Version)
}

func saveCodePackageToLocal(cds *pb.ChaincodeDeploymentSpec) error {
	if len(cds.CodePackage) == 0 {
		return errors.New("empty codepackage")
	}

	// write codepackage to local-storage
	g, err := gzip.NewReader(bytes.NewReader(cds.CodePackage))
	if err != nil {
		return fmt.Errorf("failed to init gzip reader %s", err.Error())
	}
	defer g.Close()

	baseDir := getJsChaincodeLocalPath(cds.ChaincodeSpec.ChaincodeId)
	t := tar.NewReader(g)
	for {
		hdr, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read from tar %s", err.Error())
		}
		if hdr.Typeflag != tar.TypeDir {
			fpath := filepath.Join(baseDir, hdr.Name)
			if err := os.MkdirAll(path.Dir(fpath), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create dir %s", err.Error())
			}
			ow, err := os.Create(fpath)
			if err != nil {
				return fmt.Errorf("failed to create file(%s, %s) %s", cds.ChaincodeSpec.ChaincodeId.Path, hdr.Name, err.Error())
			}

			if _, err := io.Copy(ow, t); err != nil {
				return fmt.Errorf("failed to write data %s", err.Error())
			}
			ow.Close()
		}
	}
	return nil
}

// Platform for the Javascript type
type Platform struct {
}

func (jsPlatform *Platform) ValidateSpec(spec *pb.ChaincodeSpec) error {
	return nil
}

func (jsPlatform *Platform) ValidateDeploymentSpec(spec *pb.ChaincodeDeploymentSpec) error {
	return nil
}

func (jsPlatform *Platform) GetDeploymentPayload(spec *pb.ChaincodeSpec) ([]byte, error) {

	var err error

	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tw := tar.NewWriter(gw)

	err = writeChaincodePackage(spec, tw)

	tw.Close()
	gw.Close()

	if err != nil {
		return nil, err
	}

	return inputbuf.Bytes(), nil
}

func (jsPlatform *Platform) GenerateDockerfile(cds *pb.ChaincodeDeploymentSpec) (string, error) {
	return "", nil
}

func (jsPlatform *Platform) GenerateDockerBuild(cds *pb.ChaincodeDeploymentSpec, tw *tar.Writer) error {
	if err := saveCodePackageToLocal(cds); err != nil {
		return err
	}

	return cutil.WriteBytesToPackage("codepackage.tgz", cds.CodePackage, tw)
}
