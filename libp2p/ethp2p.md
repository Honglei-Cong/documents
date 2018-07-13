## ETH P2P


eth protocol mgr管理节点的所有peers。

protocol mgr的每个peer保存对应的knownTxs 和 knownBlocks

```
ProtocolManager.BroadcastBlock(block *types.Block, propagate bool) {
	peers := pm.peers.PeersWithoutBlock(hash)
	
	// send block to subset of our peers
	transfer := peers[:Sqrt(len(peers))]
	for _, peer := range transfer {
		peer.SendNewBlock(block, td)
	}
}

```



### Discovery

func startDiscovery(context, host, gsub topicPeerLister) error

```
type topicPeerLister interface {
	ListPeers(string) []peer.ID
}

type discovery struct {
	host		host.Host
	gsub		topicPeerLister
}

func (d *discovery) HandlePeerFound(pi ps.PeerInfo)

```

### Feed

```
func (s *Server) Feed(msg interface{}) *event.Feed {
}
```

### Message

```
type Message struct {
	Peer		Peer
	Data		interface{}
}
```

### Options

```
func buildOptions() []libp2p.Option {
	// constuct example libp2p options
}
```

### Peer


### Service

```
type Server struct {
	ctx			context.Context
	cancel		context.CancelFunc
	feeds		map[reflect.Type]*event.Feed
	host		host.Host
	gsub		*floodsub.PubSub
}
```

func NewServer() (*Server, error)

1. 通过buildOptions，然后把options构建Host
2. 构建gossipSub
3. 初始化 server

func (s *Server) Start()


func (s *Server) Broadcast(msg interface{})



### Topic





