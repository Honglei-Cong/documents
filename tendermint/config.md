
## Tendermint Config


@config/config.go


```
type ConsensusConfig struct {
    RootDir                 string
    WalPath                 string
    walFile                 string

    TimeoutPropose          int           // in milliseconds
    TimeoutProposeDelta     int
    TimeoutPrevote          int
    TimeoutPrevoteDelta     int
    TimeoutPrecommit        int
    TimeoutPrecommitDelta   int
    TimeoutCommit           int

    SkipTimeoutCommit       bool          // make progress as soon as we have all the precommits

    CreateEmptyBlocks        bool
    CreateEmptyBlockInterval int

    PeerGossipSleepDuration     int
    PeerQueryMaj23SleepDuration int
}

```

