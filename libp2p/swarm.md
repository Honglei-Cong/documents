## Swarm


实现了Network接口

```
	AddTransport(t transport.Transport)
	AddAddrFilter(f string)
	Listen(addr …ma.MultiAddr)
	Process() goprocess.Process
	SetConnHandler(handler ConnHandler)
	SetStreamHandler(handler inet.StreamHandler)
	NewStreamWithPeer(ctx context.Context, p peer.ID) (*Stream, error)
	ConnectionsToPeer(p peer.ID) []*Conn
	CloseConnection(p peer.ID)
	Peers() []peer.ID
	Notify(f inet.Notifiee)
```


### Swarm

```
// 一个host的链接管理器，接收远程链接，创建远程链接，管理active的链接
type Swarm struct {
	refs			sync.WaitGroup
	local		peer.ID
	peers		pstore.Peerstore
	conns		struct {
		sync.RWMutex
		m 		map[peer.ID][]*Conn
	}
	listeners		struct {
		sync.RWMutex
		m		map[transport.Listener]struct{}
	}
	notifs		struct {
		sync.RWMutex
		m		map[inet.Notifiee]struct{}
	}
	transports 	struct {
		sync.RWMutex
		m		map[int]transport.Transport
	}

	conn		atomic.Value
	streamh		atomic.Value

	dsync		*DialSync
	backf		DialBackoff
	limiter		*dialLimiter
	Filters		*filter.Filters			// 控制哪些地址不可以链接或者被链接

	proc		goprocess.Process
	ctx			context.Context
	bwc			metrics.Reporter
}
```

* Conn：一个conn可以承载多个stream，start／close／newstream／addstream等
* Dial：
* DialPeer：连接到指定的peer.ID，返回对应的conn。Host.Connect基于此接口
* Filter：AddAddrFilter 配置地址名单
  * 在连接一个peer.ID时，查询对应的address，并过滤掉blocked address，然后选择good address
* Listen：
  * Listen开始监听一系列的地址，每个地址一个goroutine，在accept新的链接后，启动routine调用swarm.addConn，处理新的connection
* Stream：read／write
* Transport：多种链接：tcp／quic／onion等



