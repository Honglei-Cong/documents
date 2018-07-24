## PeerStore

```
type PeerStore interface {
	AddrBook				// 管理peerID与地址的对应关系，每个地址有TTL
	KeyBook				// 管理peerID与pub/priv key的关系
	Metrics					// 记录peerID与latency的对应关系
	
	Peers() []peer.ID
	PeerInfo(peer.ID) *PeerInfo
	
	Get(id peer.ID, key string) (interface{}, error)
	Put(id peer.ID, key string, val interface{}) error
	
	GetProtocols(peer.ID) ([]string, error)
	AddProtocols(peer.ID, …string) error
	SetProtocols(peer.ID, …string) error
	SupportsProtocols(peer.ID, …string) ([]string, error)
}
```



