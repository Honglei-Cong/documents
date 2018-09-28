

#### blockDAG的构建

基于UTXO，每个block需要对其中所有的交易的input utxo做依赖判断，并以此确定这个block所依赖的block，一个block可以依赖多个block。由此构建DAG。

#### identity blockchain

类似于eth的validator管理合约。

#### random generator

里面定义了一个随机数生成方式，在每个epoch保证随机分配validator到各个shard。

#### Atomix

基于Atomix协议实现cross sharding交易，Atomix基于utxo模型，由客户端协调完成跨链。具体步骤为，client将交易分别发给shardA和shardB，如果都成功了，执行成功，将结果发布到其他shard。如果有一个shard执行失败，在执行成功的shard上执行revoke操作。
1. 只能用于简单的utxo模式。
2. client端的协调不可靠。


#### shard ledger pruning

另外一个比较有意思的地方就是shard ledger pruning。
如果没有ledger pruning，validator随机分配的sharding模式就没啥意义，因为validator还是要保存所有的账本信息。
1. 每个shard的state block，类似于pbft的checkpoint。
2. 在每个epoch的结束，leader构建sharding中所有UTXO的Merkle树，计算Merkle root
3. 然后对整个Merkle root做共识，并将其保存与下一个epoch的创世块中
4. 在周期e+1 结束的时候，可以丢弃掉周期e的区块，而只保存state block和block header
5. 由client提供交易的proof of existance
6. 在周期e＋1，client可以生成周期e中的交易的proof of existence



