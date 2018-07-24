

## Tendermint Dependency Interface


#### Start／Stop／Reset／Quit
实现Service接口

#### Receive
实现Reactor接口，注册到p2p模块中
p2p在处理connection上的消息时，自动按照channel，调用对应模块的Receive接口

#### Send
调用peer的Send/TrySend接口，发送 []byte数据

#### Block Store
* Height() uint64
* LoadBlockMeta(height) *BlockMeta
* LoadBlock(height) *Block
* LoadBlockPart(height) *Part
* LoadBlockCommit(height) *Commit
* LoadSeenCommit(height) *Commit

