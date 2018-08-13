
## Quic - Go Implementation

github.com/lucas-clemente/quic-go



### Server

* quic.ListenAddr(addr, tlsConfig, config) (listener, error)
* listener.Accept() (session, error)
* session.AcceptStream() (stream, error)


### Client

* quic.DialAddr(addr, tlsConfig, config) (sessoin, error)
* session.OpenStreamSync() (stream, error)

### Session
* AcceptStream() (Stream, error)
    * returns next stream opened by peer, blocking until one is available
* AcceptUniStream() (ReceiveStream, error)
    * return the next uni-directional stream opened by peer, blocking until one is available
* OpenStream() (Stream, error)
    * opens a new bidirectional QUIC stream
    * returns a special error when the peer’s concurrent stream limit is reached
    * The peer can only accept the stream after data has been sent on the stream.
* OpenUniStream() (SendStream, error)
    * opens a new outgoing unidirectional QUIC stream
* ConnectionState() ConnectionState
    * returns basic details about the QUIC connection

### Config
* Versions []VersionNumber
    * 链接所支持的可协商的version列表，缺省支持所有verison
* RequestConnectionIDOmission
    * 请求Server端省略public header的connection ID，从而每个packet可节省8 bytes
    * 但是，将无法支持server端的IP迁移
* ConnectionIDLength
    * connectionID的字节数
* HandshakeTimeout
    * cryptographic握手所需要的时间
* IdleTimeout
    * 握手完成后，链接idle的最长时间
    * 超过idle时间后，server将自动关闭链接
* AcceptCookie func (clientAddr net.Addr, cookie *Cookie) bool
* MaxReceiveStreamFlowControlWindow
    * stream层的流量控制
    * server端缺省为1M，client端缺省为6M
* MaxReceiveConnectionFlowControlWindow
    * connection层的流量控制
    * server端缺省为1.5M，client端缺省为15M
* MaxIncomingStreams
    * 来自一个peer的stream的最大数目，缺省为100
    * 最大65535
* MaxIncomingUniStreams
    * 来自一个peer的uni stream的最大数目，缺省为100
    * 最大为65535
* KeepAlive
    * peer是否需要周期性发送PING消息，保持链接的活跃。





