
## Tendermint Peer State

PeerRoundState contains the known state of a peer.


```
type PeerRoundState struct {
    Height                     int64
    Round                      int
    Step                       RoundStepType
    StartTime                  time.Time
    Proposal                   bool
    ProposalBlockPartsHeader   PartSetHeader
    ProposalBlockParts         BitArray
    ProposalPOLRound           int
    ProposalPOL                BitArray
    Prevotes                   BitArray
    PreCommits                 BitArray
    LastCommitRound            int
    LastCommit                 BitArray
    CatchupCommitRound         int
    CatchupCommit              BitArray
}
```


