## Node

```
type Node struct {
	Identity		peer.ID

	Peerstore		pstore.Peerstore
	Reporter		metrics.Reporter
	Discovery		discovery.Service

	PeerHost		p2phost.Host
	Bootstrapper	io.Closer
	Routing			routing.IpfsRouting
	Ping			*ping.PingService			// 处理来自其它peer的ping消息

	FloodSub		*floodsub.PubSub
	PSRouter		*psrouter.PubsubValueStore
	P2P				*p2p.P2P

	proc			goprocess.Process
	ctx				context.Context
}
```



