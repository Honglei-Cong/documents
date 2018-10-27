
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



## Protocol
