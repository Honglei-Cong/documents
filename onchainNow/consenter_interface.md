


接口列表

* Consenter 
    * HandleChain(support ConsenterSupport, metadata *cb.Metadata) (Chain, error)
    	* // create and return a reference to a chain for the given set of resources.
* Chain
    * Order (env *cb.Envelope, configSeq uint64) error
    	* // Accept a message which has been processed at a given configSeq.
    	* // If the configSeq advances, consenter should revalidate the message.
    * Configure (config *cb.Envelope, configSeq uint64) error
    	* // Accept a message which reconfigure the channel and will trigger an update to the configSeq if committed.
    * WaitReady()
    	* // blocks waiting for consenter to be ready for accepting new messages.
    * Errored() <-chan struct{}
    	* // returns a channel which will close when an error has occured.
    * Start()
    	* // Start() should allocate whatever resources are needed for staying up to date with the chain.
    * Halt()
    	* // frees the resources which were allocated for this Chain.
* ConsenterSupport
    * crypto.LocalSigner
        * SignatureHeaderMaker
            * // creates a new SignatureHeader
            * NewSignatureHeader() (*cb.SignatureHeader, error)
        * Signer
            * // signer sign messages
            * Sign(message []byte) ([]byte, error)
    * msgprocessor.Processor
        * ClassifyMsg(chdr *cb.ChannelHeader) Classification
            * // classify msg to (NormalMsg, ConfigUpdateMsg, ConfigMsg)
        * ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error)
            * // validate msg based on current configuration, return current config seq
        * ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error)
            * // apply config update to the current configuration, return the resulting config message and configseq
        * ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error)
            * // ???
    * BlockCutter() blockcutter.Receiver
        * // return the block cutting helper for the channel
    * SharedConfig() channelconfig.Orderer
        * // provides the shared config from the channel's current config block
    * CreateNextBlock (messages []*cb.Envelope) *cb.Block
        * // takes a list of messages and creates the next block based on the block with highest block number comitted to the ledger
    * WriteBlock (block *cb.Block, encodedMetadataValue []byte)
        * // commits a block to the ledger
    * WriteConfigBlock (block *cb.Block, encodedMetadataValue []byte)
        * // commits a block to the ledger, and applies the config update inside
    * Sequence() uint64
        * // return the current config sequence
    * ChainID() string
        * // returns the channel ID the support is associated with
    * Height() uint64
        * // returns the number of blocks in the chain the channel is associated with





