## Sharded Byzantine Atomic Commit

* Byzantine agreement
  * ensures all honest members of a shard of size 3f+1, agree on a specific common sequence of actions
  * guarantees that all agreement is sought, a decision or sequence will eventually be agreed upon.
* Atomic commit
  * run across all shards managing objects relied upon by a transaction
  * It ensures that each shard needs to accept to be commit a transaction for the transaction to be committed; even if a single shard rejects the transaction, then all agree it is rejected.

  
#### Intial Broadcast (Prepare)

user acts as a transaction intiator, and sends 'prepare(T)' to at least one honest concerned node for transaction T.

#### Sequence Prepare

Upon a message 'prepare(T)' being received, node in each shard interpret it as the initiation of a two-phase commit protocol performed across the concerned shards.

The shard locally sequences 'prepare(T)' message through BFT protocol.

#### Process Prepare

Upon a message 'prepare(T)' being sequenced throught BFT protocol in a shard, nodes of the shard implicitly decide whether it should be committed or aborted.

Transaction T is to be committed if

* the objects input or referenced by T in the shard are active
* there is no other instance of the two-phase commit protocol on-going concerning any of those objects (no lock held)
* T is valid according to the validity rules, and the smart contract checkers in the shard

If the decision is to commit, the shard broadcasts to all concerned nodes 'prepared(T, commit)', otherwise it broadcasts 'prepared(T, abort)' -- along with sufficient signatures

The objects used or referenced by T are 'locked' in case of a 'prepared commit' until an 'accept' decision on the transaction is reached, and subsequent transactions concerning them will be aborted by the shard.

#### Process Prepared(accept or abort)

If it receives even a single 'LocalPrepared(T, abort)' from another shard, it instead will move to reach consensus on 'accept(T, abort)'.  Otherwise, if all the shards respond with 'LocalPrepared(T, commit)' it will reach a consensus with 'AllPrepared(commit, T)'.

#### Process Accept

When a shard sequences an 'accept(T, commit)' decision, it sets all objects that are inputs to the transaction T as being inactive.  It also creates any output objects from T via BFT consensus that are to be managed by the shard.  If the output objects are not managed by the shard, the shard sends requests to the concerned shareds to create the objects.


