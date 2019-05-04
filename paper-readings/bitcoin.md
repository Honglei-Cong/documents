
重读Bitcoin论文

论文题目：Bitcoin: A Peer-to-Peer Electronic Cash System

论文摘要：
Bitcoin的设计目标是一个完全去中心化的支付系统。其中重点解决的问题就是在点对点支付网络中的双花问题。
Bitcoin 通过如下方式解决双花问题：

1. 链式区块结构，所有交易保持在链式结构中
2. 基于hash方式的PoW
3. 基于最长链的容错方式


引言部分：

首先定义目前支付系统的问题：trust based model，无法实现不可撤销的支付，也就无法完全避免双花。
而且中间人带来的交易摩擦费也增加了小额支付的成本。

'With the possibility of reversal, the need for trust spreads.'
对信任的需求，对交易的不确定性，都是交易成本提高的因素。以此，Bitcoin就是在支付领域解决这个问题。


因此，Bitcoin设计了链式的交易结构。
miner保持网络中所有的交易，因此可以独立验证每个区块中交易的合法性。
然后，通过PoW保证所有miner的合法性。


PoW的设计思想

'Once the CPU effort has been expended to make it satisfy the proof-of-work, the block cannot be changed without redoing the work.'
'As later blocks are chained after it, the work to change the block would include redoing all the blocks after it.'
'The proof-of-work difficulty is determined by a moving average targeting an average number of blocks per hour.'

PoW + Longest chain，一个全新的共识算法的诞生。



Bitcoin对于网络运行环境的设计：

* New transactions are broadcast to all nodes.
* When a node finds a proof-of-work, it broadcasts the block to all nodes.
* Nodes express their acceptance of the block by working on creating the next block in the chain, using the hash of the accepted block as the previous hash.


Bitcoin的激励设计：
与PoW一起开创了加密经济的先河

* The first transaction is a special transaction that starts a new coin owned by the creator of the block
* Once a predetermined number of coins have entered circulation, the incentive can transition entirely to transaction fees and be completely inflation free.

Bitcoin是完全0通胀的token设计，这是Bitcoin设计中很重要的一项，但是后续很多PoW没有延续这个设计。
Bitcoin初始就是作为一种支付工具进行的设计，而不是作为一种价值承载工具。
而且，在Bitcoin初始设计中也不认为token丢失或者会存在屯币的市场行为，因为Bitcoin只是被设计为一个支付的媒介，并没有考虑如果全世界都使用Bitcoin进行支付会怎样。

Bitcoin的发展应该是超出了设计者最初的想象。


Bitcoin的存储设计

Bitcoin的存储设计中介绍了每个矿工都需要维护哪些数据和如何维护这些数据。由于每个矿工都是全量节点，也就是Bitcoin整个网络的数据量估计。

在Bitcoin的初始设计中充分考虑了随着时间增长给矿工节点带来的存储压力，Bitcoin的UTXO设计使得矿工节点对于每个区块只需要保存少量的数据即可完成链上所有交易的验证。
而更多是是由交易的发起者提供交易中UTXO的有效证明。
这个设计和传统的客户端／服务端设计完全不同，而是很好的延续了Bitcoin的去中心化设计思想，矿工节点的责任是验证交易，而是用户的责任是提供有效交易证明。
在Bitcoin网络中，用户和矿工只是角色的不同，而不是服务与被服务的关系，Bitcoin的每个用户维护自己的UTXO，矿工只为Bitcoin网络激励而服务。

但是在当前的Bitcoin生态中，由于区块浏览器和各种Wallet应用的普及，大部分用户也只保存自己的私钥和地址，而很少运行自己的SPV节点。

与Bitcoin存储设计相辅相成的是其UTXO的交易模式。
[TODO]



Bitcoin的隐私保护
[TODO]

