## Stream


```
type Stream peerstream.Stream

	InGroup(g Group) bool
	AddGroup(g Group) error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	SetDeadline(t time.Time)
```

### peerstream.Stream

```
	type Stream struct {
		smuxStream 	smux.Stream			// stream-muxer.Stream
		conn 			*Conn
		groups 		groupSet
		protocol 		protocol.ID
	}

	type stream-muxer.Stream interface {
		io.Reader
		io.Writer
		io.Closer
		Reset() error
		SetDeadline(time.Time) error
		SetReadDeadline(time.Time) error
		SetWriteDeadline(time.Time) error
	}

	type stream-muxer.Conn interface {
		io.Closer
		IsClosed() bool
		OpenStream() (stream, error)
		AcceptStream() (stream, error)
	}

	type stream-muxer.Transport interface {
		NewConn(c net.Conn, isServer bool) (Conn, error)
	}
```


