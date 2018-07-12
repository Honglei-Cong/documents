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



