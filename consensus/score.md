
# Score

A scalable one-copy serializable partial replication protocol


Two scalar variables per node:

* commitId
    * maintains the timestamp that was attributed to the last update transaction when committed on that node.
* nextId
    * keep tracks of the next timestamp that the node will propose when it will a commit request for a transaction that accessed some of the data that it maintains.

## Snapshot visibility

Snapshot Visibility for transactions is determined by associating with each transaction T a scalar timestamp, which we call *snapshot identifier* or, sid.

The sid of a transaction is established upon its first read operation.  

In this case the most recent version of the requested datum is returned, and the transaction's sid is set to the value of the commitId at the transaction's originating node, if the read can be served locally.  Otherwise, if the requested datum is not maintained locally, T.sid is set equal to the maximum between *commitId* at the originating node and *commitId* at the remote node from which T reads.

From that moment on, any subsequent read operation is allowed to observe the most recent committed version of the requested datum having timestamp less thatn or equal to T.sid, as in classical multiversion concurrency control algorithms.


## Atomic commit

A key mechanism used in SCORe to correctly serialize transactions, and in particular to track write-after-read dependencies, is to update the *nextId* of a node upon processing of a read operation.

If a node receives a read operation from a transaction T having a *sid* larger than its local *nextId*, this is advanced to T.sid.

Since a transaction is attributed a snapshot identifier upon its first read, which is used throughout its execution, SCORe guarantees that the snapshot read by a transaction is alwasy consistent with respect to a prefix of the history of committed transactions.

As a consequence, in SCORe read-only transactions never abort and do not need to undergo any distributed verification.


