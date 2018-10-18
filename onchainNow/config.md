
一个consenter的配置分为两部分

* metadata config, 保存于block
* shared config, 保存于共享配置中


### Metadata in block

保存每个区块需要的信息

Metadata index in block

* BlockMetadataIndex_SIGANTURES		: for block signatures
* BlockMetadataIndex_LAST_CONFIG		: last config block seqno
* BlockMetadataIndex_TRANSACTION_FILTER    : bit array filter for invalid transactions
* BlockMetadataIndex_ORDERER			: operational metadata for orderers


### Shared Config

(from tendermint)

Base Config

* chainID string
* RootDir string
* FastSync bool

P2P Config

* ListenAddress string
* Seeds string
* PersistentPeers string
* HandshakeTimeout time.Duration
* DailTimeout time.Duration

Consensus Config

* WalPath string
* WalFile string
* Timeout(Propose, Prevote, Precommit, Commit) int
* CreateEmptyBlocks bool
* CreateEmptyBlocksInterval int
* PeerGossipSleepDuration int
* PeerQueryMaj23SleepDuration int

I
