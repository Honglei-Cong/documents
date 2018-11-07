
# Non-Monotonic Snapshot Isolation

* RC: read committed
* NMSI: non-monotonic snapshot isolation
* PSI: parallel snapshot isolation
* US: update serialisability
* SER: serialisability

NMSI: 

* both satisfies strong safety properties
* and addresses the scalability problems of PSI

SI: a transaction reads its own consistent snapshot, and aborts iff its writes conflicts with a previously-committed concurrent transaction.

### Contribution

Identified the following essential scalability properties:

* only replicas updated by a transaction T make steps to execute T
* a read-only transaction never waits for concurrent transactions and always commit
* a transaction may read object versions committed after it started
* two transactions synchronise with each other only if their writes conflict

### Scalability Properties

* Wait-Free Queries
    * Read-only transactions does not wait for concurrent transactions and always commits.  
    * (For blockchain, read-only transactions are executed off-chain.)
* Genuine Partial Replication
    * Many partial replication protocols perform system wide global consensus or communication.
    * Genuine Partial Replication, transactions communicates only with replicas that store some object accessed in the transaction.
    * Non-conflicting transactions do not interfere with each other, and the intrinsic parallelism of a workload can be exploited.
    * (good, only related smart-contracts are locked in current design.)
* Minimal Commitment Synchronization
    * To keep consensus, while alleviating their costs
    * during commitment, transaction T_i waits for transaction T_j only if T_i and T_j write-conflict. 
    * (What is write-conflict? seems read-write conflict need to be handled in current design)
* Forward Freshness
    * A criterion support Forward Freshness if it allows reading an object version that committed after the start of the transaction.
    * In case of global transactions (i.e., transactions that touch several sites), this property is fundermental.


## Definition of NMSI

**Dependency**

Consider a history h and two transactions T_i and T_j.

We note T_i » T_j when the transaction T_i reads a version of x installed by T_i (i.e. r_i(x_j) is in h)

Transaction T_i *depends* on transaction T_j when the above relation holds by transitivity, that is T_i »* T_j.

Transaction T_i and T_j are *independent* if either T_i »* T_j nor T_j »* T_i


**Consistent Snapshot (CONS)**

(A transaction sees a consistent snapshot iff it observes the effects of all transactions it depends on.)

Ti in a history h observes a consistent snapshot iff, for every object x, if T_i reads version x_j, T_k writes version x_k, and T_i depends on T_k, then version x_k is followed by version x_j in the version order included by h (i.e. x_k «_h x_i).

We write h is CONS when all transactions in h observe a consistent snapshot.

**Avoiding Cascading Aborts (ACA)**

(prevent transactions to read non-committed data)

A history h avoids cascading aborts when for every read r_i(x_j) in h, operation c_j precedes r_i(x_j) in h.

ACA denotes the set of histories that avoid cascading aborts.

**Write-Conflict Freedom (WCF)**

(forbids independent write-conflicting updates to commit)

A history h is write-conflict free, iff independent committed transactions never write to the same object.

**NMSI**

A history h is in NMSI iff h belongs to (ACA & CONS & WCF)


* 如果 T_i 读取了 T_j 写入的数据，T_i 依赖于 T_j
* 依赖关系具有传递性
* 所有事务都从一致性快照中读取数据，比如 T_i 读取到了 x_j，而 T_k 写入了 x_k，T_i 依赖于 T_k，那么 x_k 的写入肯定发生在 x_j之前。
* 所有事务读取的都是commit后的数据。 数据只有commit后才能被其它事务读取。
* 相互独立的事务，不会对写入到同一数据。即，如果两个事务写入到了同一数据，他们之间必定存在依赖关系。

通过rwlock实现ACA，可以提高跨链效率。


## Protocol

### state of Transaction:

* Executing
  * Each non-termination operation o_i(x) in T_i is executed optimistically at the transaction coordinator coord(T_i).
  * if o_i is read, coord(T_i) returns the value, fetched either from local replica or a remote one.
  * if o_i is write, coord(T_i) stores the corresponding update valud in a local buffer, enabling
      * subsequent reads to observe the modification
      * subsequent commit to send the write-set to remote replicas
* Submitted
  * Once all read/write operations of T_i have executed, T_i terminates, and the coordinator submits it to the termination protocol.
  * The protocol applies a certification test on T_i to enforce NMSI.
  * This test ensures that if two concurrent conflicting update transactions terminate, one of them aborts.
* Committed/Aborted
  * If Committed, its updates are applied to the local data store.
  * If Aborted, T_i is aborted.


### Termination Protocol

In order to satisfy GPR, the termination protocols uses a genuine atomic multicast primitives.  This requires:

* we form non-intersecting groups of replicas, and an eventually leader oracle is available in each group
* a system-wide reliable failure detector is available.


