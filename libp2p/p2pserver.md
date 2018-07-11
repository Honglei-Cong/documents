# P2P Server

## function list

* initialization


* start

1. network start
2. msg router start
3. connect recent peers
4. connect seed nodes
5. sync with recent peers
6. start timer to keep sync with inactive peers
7. start heartbeat service
8. start block syncing


* stop


* block sync




## server state

```
type Server {
	* P2P network
	* message router
	* block sync manager
	* recent peers
}
```

msg router的RecvChan接收网络上的所有payload，
对每一种payload，找到对应的msg handler，启动go routine处理消息。

```
type MessageRouter {
	msgHandlers		map[string]MessageHandler
	RecvChan			chan *MsgPayload
	stopCh				chan bool
}
```

```
type NetServer struct {
	* self peer info
	* listener
	* msgPayload channel
	* peer address map
	* neighbours
}
```

```
type PeerInfo struct {
	* id
	* version
	* port
	* block height
}
```

```
type Link struct {
	* remote peer id
	* remote peer addr
	* remote peer port
	* connection net.Conn
	* recvChan chan Payload
}
```


## Actor

actor处理的消息

* StartReq
* StopReq
* GetPortReq
* GetVersionReq
* GetConnectionCntReq
* GetPortReq
* GetNeighbourAddressesReq
* AppendPeerID
* RemotePeerID
* AppendHeaders
* AppendBlocks
* TransmitConsensusMsgReq


## Msg Routing

根据msg header，将msg做反序列化，然后根据消息类型处理消息

然后通过eventbus通知到对应的模块

```
type messageHeader struct {
	Magic		uint32
	CMD			[]byte
	Length		uint32
	Checksum	[]byte
}
```


## Block Syncer




