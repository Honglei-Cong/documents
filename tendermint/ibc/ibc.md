# Cosmos IBC Protocol

IBC = Inter-Blockchain Communication

IBC是cosmos定义的链间通讯的一个协议，定义了一个 stricted-ordered messaage passing 通道。和我们当前的shard设计中的链间通讯部门有很多相同之处，同时我们也可以从Cosmos IBC中借鉴一些好的细节上的设计。

与 shard 设计不同，cosmos只为链间的资产交换而设计，而sharding目标是链间的事务交互，所以可借鉴的部分仅限于 strict-ordered message passing channel 这个部分。

IBC 要求每个子链:

* cheaply verifiable rapid finality
* Merkle tree substate proof


## IBC 中定义的一些概念

```
Packet = (type, sequence, source, destination, data)
Receipt = (sequence, source, destination, result)
Queue = (q_head, index, q_tail)
```

**Ordered Queue**

* can be conceptualized as a slice of an infinite array.
* elements can be appended to the tail and removed from the head.
* *advance* to facilitate efficient queue cleanup

**Channel**

a set of the required packet queues to facilitate ordered bidirectional communication between two blockchains A and B.


## Requirements

每个队列都是以Merkle树的形式保存，队列中每个消息的Key定义为：

```
key = (queue name, head | tail | index)
```

每个队列的实现，和blockchain类似，一旦加入到了队列中，将不可更改。

每个队列都是在可验证的chain间建立，即在将一个元素放入队列前，必须检查此元素的确是来自对方队列中。


### IBC Queue Format

```
key = (remote id, [send | receipt], [head | tail | index])
V_send = (maxHeight, maxTime, type, data)
V_receipt = (result, [sucess | error code])
```

key和value都将保存在Merkle树中，Merkle树的root hash将被共识保存在区块头中。
其它的chain通过spv模式，验证消息的有效性。


