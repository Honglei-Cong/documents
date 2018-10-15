
## Tendermint Config


@config/config.go


```
// Config defines the top level configuration for a Tendermint node
type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for services
	RPC             *RPCConfig             `mapstructure:"rpc"`
	P2P             *P2PConfig             `mapstructure:"p2p"`
	Mempool         *MempoolConfig         `mapstructure:"mempool"`
	Consensus       *ConsensusConfig       `mapstructure:"consensus"`
	TxIndex         *TxIndexConfig         `mapstructure:"tx_index"`
	Instrumentation *InstrumentationConfig `mapstructure:"instrumentation"`
}



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

