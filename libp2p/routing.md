## Routing

```
type IpfsRouting interface {
	ContentRouting
	PeerRouting
	ValueStore

	// 通知路由系统进入bootstrap状态
	Bootstrap(context) error
}

// 用于寻找那个peer拥有哪些信息
type ContentRouting interface {
	Provide(context, *cid.Cid, bool) error
	FindProviderAsync(context, *cid.Cid, int) <- chan pstore.PeerInfo
}
type Cid struct {			// 信息的哈希
	version uint64
	codec uint64
	hash mh.Multihash
}
type PeerInfo struct {		// 记录每个peer的ID和peer的地址
	ID peer.ID
	Addrs []Multiaddr
}

// 用于基于peer ID寻找peer的信息
type PeerRouting interface {
	FindPeer(context, peer.ID) (pstore.PeerInfo, error)
}

// 提供数据的Get／Put接口
type ValueStore interface {
	PutValue(context, string, []byte, …Options) error
	GetValue(context, string, …Options) ([]byte, error)
}
```



