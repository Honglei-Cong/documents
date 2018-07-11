## Host

```
type Host interface {
	ID()				peer.ID
	Peerstore()		pstore.Peerstore

	// host的所有listen的地址
	Addrs()			[]ma.Multiaddr

	// host的network interface
	Network()		inet.Network
	Mux()			*msmux.MultistreamMuxer
	Connect(context, pi pstore.PeerInfo) error

	// 设定Mux的消息处理函数
	SetStreamHandler(pid protocol.ID, handler inet.StreamHandler)

	// 基于match函数的mux消息处理handler
	SetStreamHandlerMatch(protocol.ID, func(string) bool, inet.StreamHandler)
	RemoveStreamHandler(pid protocol.ID)
	NewStream(context, p peer.ID, pids …protocol.ID) (inet.Stream, error)
	Close() error
	ConnManager() ifconnmgr.ConnManager
}
```

### BlankHost

```
type BlankHost struct {
	n		inet.Network
	mux	*mstream.MultistreamMuxer
	cmgr	ifconnmgr.ConnManager
}
```



