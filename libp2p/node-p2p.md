## Node P2P

```
type P2P struct {
	Listeners		ListenerRegistry
	Streams			StreamRegistry
	identity			peer.ID
	peerHost		p2phost.Host
	peerstore		pstore.Peerstore
}
type Listener interface {
	Accept() (net.Stream, error)
	Close() error
}
type ListenerRegistry struct {
	Listeners 		[]*ListenerInfo
}
type StreamRegistry struct {
	Streams 		[]*StreamInfo
	nextID 			uint64
}
```



```
func (p2p *P2P) newStreamTo (ctx, peer.ID, protocol string) (net.Stream, error)
func (p2p *P2P) Dial (ctx, ma.Multiaddr, peer.ID, proto string, bindAddr Multiaddr) (*ListenerInfo, error)
func (p2p *P2P) doAccept (*ListenerInfo, remote net.Stream, listener manet.Listener)


func (p2p *P2P) registerStreamHandler(ctx, protocol) (*P2PListener, error)
func (p2p *P2P) NewListener(ctx, protocol string, addr Multiaddr) (*ListenerInfo, error)
func (p2p *P2P) acceptStreams(*ListenerInfo, Listener)
func (p2p *P2P) CheckProtoExists(proto string) bool


```