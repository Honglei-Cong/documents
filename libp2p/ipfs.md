## IPFS on libp2p



* Host
  * constructPeerHost

* peer store @ go-libp2p-peerstore



### Config

```
type Config struct {
	PeerKey			crypto.PrivKey
	
	Transports				[]TptC
	Muxers						[]MsMuxC
	SecurityTransports		[]MsSecC
	Insecure					bool
	Protector					pnet.Protector
	
	Relay				bool
	ListenAddrs		[]Multiaddr
	Filters			*filter.Filters
	
	NATManager		NATManagerC
	Peerstore			pstore.PeerStore
	Reporter			metrics.Reporter
}
```

#### Config init

Cfg.Apply(opts ...Option) error

* constructDHTRouting
* constructPeerHost


### Network (Swarm)

@go-libp2p-swarm

swarm.NewSwarm(ctx, pid, peerstore, reporter)


### BasicHost

```
type struct BasicHost {
	network			inet.Network
	mux					*msmux.MultistreamMuxer
	ids					*identity.IDService
	natmgr				NATManager
	addrs				AddrsFactory
	
	maResolver		*madns.Resolver
	cmgr				ifconnmgr.ConnManager
	
	proc				goprocess.Process
}
```

### Msg Format

```
bitswap msg

@exchange/bitswap/message/message.go
type impl struct {
	full			bool
	wantlist		map[string]*Entry
	blocks			map[string]blocks.Block
}

```

message提供到protobuf的转换方法，stream writer通过protobuf标准接口将msg序列化，然后发送。


### Send Message

exchange/bitswap/network/ipfs_impl.go

streamMessageSender::SendMsg(msg bsmsg.BitSwapMessage) error


bitswap的msg queue初始化时，指定对应的peer，所以在NewMessageSender中创建stream到peer。

wantManager.startPeerHandler，每个新的peer的connectEvent，wangManager启动对应peer的msg queue。

connectEvent来自netNotifiee。


通过ValueStore的接口使用了DHT，通过findProvider接口查询对应的peer info。
unicast是通过messageSender创建stream，然后在stream上发送。

也就是libp2p提供的是基于内容的DHT网络。对于blockchain需要的是基于peer info的DHT网络？

IpFSDHT.strmap保存了peerID到msg sender的mapping。



