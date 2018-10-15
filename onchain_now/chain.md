Chain需要实现如下接口：

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


### Start

1. 通过 ConsenterSupport 加载本地配置
2. 完成所有相关模块的初始化
3. 启动 Start Thread


### WaitReady

接收到一个broadcast的grpc请求，在处理过程中等待processor就绪。

WaitReady完成后，调用ProcessNormalMsg/Order 或者 ProcessConfigUpdateMsg/Configure 处理请求。


### Order

To put tx to mem queue.

When constructing new block, using block cutter to cut block from mem queue.

On block consensused, remove txs from pool.  New blocks are saved with ConsenterSupport.WriteBlock.

ConsenterSupport will broadcast all blocks to peers when WriteBlock done.



