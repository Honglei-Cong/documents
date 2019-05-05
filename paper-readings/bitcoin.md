
重读Bitcoin论文

论文题目：Bitcoin: A Peer-to-Peer Electronic Cash System

论文摘要：
Bitcoin的设计目标是一个完全去中心化的支付系统。其中重点解决的问题就是在点对点支付网络中的双花问题。
Bitcoin 通过如下方式解决双花问题：

1. 链式区块结构，所有交易保持在链式结构中
2. 基于哈希计算方式的工作量证明PoW
3. 基于最长链的容错方式


引言部分：

首先定义目前支付系统的问题：trust based model，无法实现不可撤销的支付，也就无法完全避免双花。
而且中间人带来的交易摩擦费也增加了小额支付的成本。

'With the possibility of reversal, the need for trust spreads.'
对信任的需求，对交易的不确定性，都是交易成本提高的因素。以此，Bitcoin就是在支付领域解决这个问题。


因此，Bitcoin设计了链式的交易结构。
所有miner都是平等的，每个miner都保存网络中所有的交易的执行结果，可以独立验证每个区块和区块中每个交易的合法性。
然后，通过PoW + Longest Chain的共识算法协调所有miner共同维护Bitcoin网络协议的运行。

## PoW + Incentive

Bitcoin对于网络运行环境的设计：

* New transactions are broadcast to all nodes.
* When a node finds a proof-of-work, it broadcasts the block to all nodes.
* Nodes express their acceptance of the block by working on creating the next block in the chain, using the hash of the accepted block as the previous hash.

Bitcoin运行于完全公开的网络环境中，所有人可以随时加入到网络中，网络中所有信息都是完全公开的方式。
由于Bitcoin的完全公开完全平等的网络环境，如何使Bitcoin协议能够长久稳定运行，就是Bitcoin设计最杰出的地方，同时也是加密经济系统设计方法的创建点。

首先是Bitcoin的PoW设计思想

'Once the CPU effort has been expended to make it satisfy the proof-of-work, the block cannot be changed without redoing the work.'
'As later blocks are chained after it, the work to change the block would include redoing all the blocks after it.'
'The proof-of-work difficulty is determined by a moving average targeting an average number of blocks per hour.'

在一个公开匿名且没有监管网络环境中，如何吸引矿工参与到Bitcoin网络，并如何协调所有矿工运行Bitcoin协议。
首先第一个问题就是，如何判断一个矿工参与到了Bitcoin网络中，即矿工如何证明自己对Bitcoin网络所做的贡献。
这个问题就是通过PoW的方式实现。
PoW算法建立在密码学哈希算法的随机性假设的基础上，Bitcoin采用SHA256哈希算法，即原数据的任意改动通过SHA256哈希后都将产生完全随机的另一个结果。
在Bitcoin协议中定义了区块的困难度，使得每个区块的构建都附加着完成一定困难度的计算量。
密码学哈希算法的完全随机性为Bitcoin带来的极大的安全性，由于Bitcoin区块的链式结构，对任意历史区块的修改都必须要付出从那个区块高度到最新区块高度所有的计算量的总和。
同时完全随机性也是一种公平性，从而实现在没有监管的网络中的公平竞争。

因此，PoW + Longest chain规则，形成了一个全新的共识算法，它首次解决了匿名无监管网络中的一致性问题。

PoW共识算法定义了参与规则后，Bitcoin的激励设计保证了Bitcoin网络的持续运行。
Bitcoin正是通过激励设计＋PoW共识算法，开创了加密经济的先河。

## Bitcoin的经济模型

* The first transaction is a special transaction that starts a new coin owned by the creator of the block
* Once a predetermined number of coins have entered circulation, the incentive can transition entirely to transaction fees and be completely inflation free.

Bitcoin是完全0通胀的token设计，这是Bitcoin设计中很重要的一项，但是后续很多PoW没有延续这个设计。
Bitcoin初始就是作为一种支付工具进行的设计，而不是作为一种价值承载工具。
而且，在Bitcoin初始设计中也不认为token丢失或者会存在屯币的市场行为，因为Bitcoin只是被设计为一个支付的媒介，并没有考虑如果全世界都使用Bitcoin进行支付会怎样。

Bitcoin的发展应该是超出了设计者最初的想象。


## Bitcoin的存储设计

Bitcoin的存储设计中介绍了每个矿工都需要维护哪些数据和如何维护这些数据。由于每个矿工都是全量节点，也就是Bitcoin整个网络的数据量估计。

在Bitcoin的初始设计中充分考虑了随着时间增长给矿工节点带来的存储压力，Bitcoin的UTXO设计使得矿工节点对于每个区块只需要保存少量的数据即可完成链上所有交易的验证。
而更多是是由交易的发起者提供交易中UTXO的有效证明。
这个设计和传统的客户端／服务端设计完全不同，而是很好的延续了Bitcoin的去中心化设计思想，矿工节点的责任是验证交易，而是用户的责任是提供有效交易证明。
在Bitcoin网络中，用户和矿工只是角色的不同，而不是服务与被服务的关系，Bitcoin的每个用户维护自己的UTXO，矿工只为Bitcoin网络激励而服务。

但是在当前的Bitcoin生态中，由于区块浏览器和各种Wallet应用的普及，大部分用户也只保存自己的私钥和地址，而很少运行自己的SPV节点。

与Bitcoin存储设计相辅相成的是其UTXO的交易模式。
基于UTXO的支付方式和互联网中的电子支付完全不同，而更接近于现实中的实物支付。
UTXO的支付交易中包括输入和输出两个部分，输入部分的Bitcoin将被花费掉，用于生产出输出部分的未花费Bitcoin。
可以看到，Bitcoin的交易处理过程更像是一个实物生产过程，而合电子支付中的余额转账方式完全不同。

通常输入部分的Bitcoin总和等于生产出来的未花费Bitcoin总和。如果输入的总和大于输出的总和，多余部分将返还交易的发送者。但是如果输入的总和小于输出的总和，这将是一个无效的交易。这个过程可以类比于工厂购买原材料加工生产出新的产品的过程，虽然所有交易都会产生未花费的Bitcoin，但是不同交易产出的Bitcoin是有不同的编号，类似于同一类产品的不同编号。

正是这种类似工厂生产的UTXO交易模型，使得Bitcoin矿工可以只维护当前未花费Bitcoin状态，即可对所有当前交易进行有效验证。


## Bitcoin的隐私保护

Bitcoin网络中所有交易都在网络中公开广播，一个用户如果要基于Bitcoin做支付，那么他必须公开自己发出的所有交易。
这一点使得Bitcoin看起来没有任何隐私保护。但是，Bitcoin网络中所有信息都是和公私钥绑定的，而无法从公私钥信息在Bitcoin网络中直接反向到真实的用户。
而且，用户还可以通过不同交易中采用不同的公私钥的方式可以很好地保护自己的隐私。

