## Network

```
type Network interface {
	Dialer
	io.Closer
	SetStreamHandler(StreamHandler)
	SetConnHandler(ConnHandler)

	// 返回和peer.ID的链接，如果不存在将新创建一个
	NewStream(context, peer.ID) (Stream, error)
	Listen(…ma.Multiaddr) error

	// 返回network正在listen的所有地址
	ListenAddresses() []Multiaddr
	InterfaceListenAddresses() ([]MultiAddr, error)

	// network对应的process
	Process() goprocess.Process
}
```

swarm实现了Network接口

