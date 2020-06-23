
在执行区块中交易时，无需访问完整world state。

区块创建者提取执行交易所需要的所有内容（block witness），打包到附属数据结构，和区块一起在网络中传播。
基于附属数据，客户端可以重建出执行区块中交易所需要的所有信息，从而验证区块中的交易。

block witness，包括每个交易的read set和write set所对应的merkle tree。

按照对之前ethereum的统计，每个区块需要的witness数据在1MB左右。

* 将ethereum的MerkleTree由16分叉改为二叉树，使得witness中的哈希数量大大减少，可以将witness数据大小降低一半。
* 将代码分块，减小代码。

### Polynomial Commitment

通过多项式承诺将witness大小减小到几十KB。



