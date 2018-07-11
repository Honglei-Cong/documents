## Transport

```
type Transport interface {
	Dialer(local_addr, opt) (Dialer, error)
	Listen(local_addr) (Listener, error)
	Matches(addr) bool
}

type Dialer interface {
	Dial (remote_addr) (Conn, error)
	DialContext (context, remote_addr) (Conn, error)
	Matches(addr) bool
}

type Listener interface {
	Accept() (Conn, error)
	Close() error
	Addr() net.Addr
	Multiaddr() ma.Multiaddr 
}
```

