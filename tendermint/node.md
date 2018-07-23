
## Tendermint Node State

@consensus/state.go


```
type ConsensusState struct {
    cmn.BaseService

    config                *cfg.ConsensusConfig

    blockStore        sm.BlockStore
    mempool          sm.Mempool            // 交易池
    evpool              sm.EvidencePool

    RoundState

    peerMsgQueue            chan msgInfo
    internalMsgQueue        chan msgInfo
    timeoutTicker                TimeoutTicker

    wal                                WAL
    replayMode                  bool
    doWALCatchup           bool
}
```

#### 数据通道

* StateChannel
    * NewRoundStepMessage
    * CommitStepMessage
    * HasVoteMessage
    * VoteSetMaj23Message
    * ProposalHeartbeatMessage
* DataChannel
    * ProposalMessage
    * ProposalPOLMessage
    * BlockPartMessage
* VoteChannel
    * VoteMessage
* VoteSetBitsChannel
    * VoteSetBitsMessage



