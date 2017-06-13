
# Practical Byzantine Fault Tolerance

The algorithm offers both liveness and safety provided at most lower(n-1/3) out of a total of n replicas are simultaneously faulty.

It uses only one message round trip to execute read-only operations and two to execute read-write operations.  Also, it uses an efficient authentication scheme based on message authentication codes during normal operation: public-key cryptography is used only when there are faults.


## system model

Byzantine failure model: faulty nodes may behave arbitrarily, subject only to the below restriction:
>
1. we assume independent node failures (N-version programming)
2. use cryptography tech to preventing spoofing and replays and to detect corrupted messages.
3. our messages contains public-key signatures, message authentication codes, and message digest produced by collision-resistant hash functions
4. allow for a very strong adversary that can coordinate faulty nodes, delay communications, or delay correct nodes in order to cuase the most damage to the replicated service.  But, we do assume that th adversary cannot delay correct nodes indefinitely.
5. we also assume that the adversary are computationally bound so that it is unable to subvert the cryptography tech.

## Service Properties

**safety**: means that the replicated service satisfies linearizability; it behaves like a centralized implementation that executes operations atomically one at a time.  

Safety is provided regardless of how many faulty clients are using the service; all operations performed by faulty clients are observed in a consistent way by non-faulty clients.  In particular, if the service operations are designed to preserve some invariants on the service state, faulty clients cannot break those invariants.

**Liveness**: The algorithm does not rely on synchrony to provide safety.  Therefore, it must rely on synchrony to provide liveness.  We guarantee liveness, i.e., clients eventually receive replies to their requests, provided at most lower(n-1/3) replicas are faulty, and delay(t) does not grow faster than t indefinitely.  delay(t) is the time between the moment t when a message is sent for the first time and the moment when it is received by its destination.

**Fault-tolerant privacy**:  The algorithm does not address the problem of fault-tolerant privacy.

## The Algorithm

Our algorithm is a form of state machine replication: the service is modeled as a state machine that is replicated across different nodes in a distributed system.

The replicas move through a succession of configurations called **views**.  In a view, one replica is the **primary**, and the others are **backups**.  View are numbered consecutively.  view changes are carried out when it appears that the primary has failed.

The algorithm works roughly as follows:
>
1. A client sends a request to invoke a service operation to the primary
2. The primary multicasts the request to the backups.
3. replicas execute the request and send a reply to the client
4. the client wait for f+1 replies from different replicas with the same result; this is the result of the operation.

two requirements on replicas:
>
1. they must be deterministic
2. they must start in the same state

### The Client

client requests the execution of state machine operation o by sending request:
>
1. tag-request
2. operation
3. timestamp, to ensure exactly-once semantic
4. client-id

Each message sent by the replicas to the client includes the current view number, allowing the client to track the view and hence the current primary.

A client sends a request to what it believes is the current primary using a point-to-point message.  The primary atomically multicasts the request to all the backups using the protocol described later.

A replica sends the reply to the request directly to the client.  The repliy contains:
>
1. tag-reply
2. current view number
3. the timestamp of the corresponding request
4. client-id
5. replica-number
6. result of executing the requested operation

If the client does not receive replies soon enough, it broadcasts the request to all replica.  If the request has been already been processed.


### Normal-Case Operation

state of each replica includes:
>
1. state of the service
2. a message log containing messages the replica has accepted
3. an integer denoting the replica's current view

When the primary, p, receives a client request, m, it starts a three-phase protocol to atomically multicast the request to replica: pre-prepare, prepare, commit

The pre-prepare and prepare phases are used to totally order requests sent in the same view even when the primary, which proposes the ordering of requests, is faulty.

The prepare and commit phases are used to ensure that requests that commit are totally ordered across views.

***Pre-Prepare Phase***
>
1. the primary assigns a sequence number, n, to the request
2. multicast a pre-prepare message with m piggybacked to all the backups
3. appends the message to its log

the message contains:
>
1. tag-pre-prepare
2. view number
3. sequence number
4. digest of client request

Requests are not included in pre-prepare messages to keep them small.  This is important because pre-prepare messages are used as a proof that the request was assigned sequence number n in view v in view changes.  Additionally, it decouples the protocol to totally order requests from the protocol to transmit the request to the replicas.

A backup accepts a pre-prepare message provided:
>
1. signatures in the request and the pre-prepare message are corrent and d is digest for m
2. it is in view v
3. it has not accepted a pre-prepare message for view v and sequence number n containing a different digest
4. the sequence number in the pre-prepare message is between a lower water mark, h, and a high water mark H.

***Prepare Phase***

If backup i accepts the PRE-PREPARE message, it enters the prepare phase by **multicasting** message, which contains:
>
1. tag-prepare
2. view number
3. sequence number
4. digest of client request
5. self-node id


If not accept the PRE-PREPARE message, do nothing.


A replica (including the primary) accepts prepare messages and adds them to its log provided their signatures are correct, their view number equals the replica's current view, and their sequence number is between h and H.

define predicate **prepared(m, v, n, i)** to be true if and only if: 
> 
replica i has inserted in its log, the request m, a pre-prepare for m in view v with sequence number n, and (2 x f) nodes prepares from different backups that match the pre-prepare.  The replica verify whether the prepares match the pre-prepare by checking that they have the same view, sequence number, and digest.

The pre-prepare and prepare phases of the algorithm guarantee that non-faulty replicas agree on a total order for the request within a view.


replica i multicasts a COMMIT message to the other replicas, when prepared(m, v, n, i) becomes true.  COMMIT message contains:
>
1. tag-COMMIT
2. v, view-number
3. n, sequence-number
4. D(m), digest of client request
5. i, self-node id

***Commit Phase***

This starts the commit phase.  Replicas accept commit messages and insert them in their log provided they are properly signed, the view number in the message is equal to the replica's current view, and the sequence number is between h and H.

define the predicate **committed(m, v, n)**, is true if and only if 
>
1. prepared(m, v, n, i) is true for all i in some set of f+1 non-faulty replicas
2. commited-local(m, v, n, i) is true if and only if prepared(m, v, n, i) is true and i has accepted 2f+1 commits from different replicas that match the pre-prepare for m:  A commit matches a pre-prepare if they have the same view, sequence number, and digest.

commit phase ensures the following invariant: if committed-local(m, v, n, i) is true for some non-faulty i, then committed(m, v, n) is true.

Each replica i executes the operation requested by m after committed-local(m, v, n, i) is true, and i's state reflects the sequential execution of all requests with lower sequence number.  After executing the requested operation, replicas send a reply to the client.


### Garbage Collection

checkpoint and stable checkpoint

A replica maintains several logical copies of the service state:
>
1. the last stable checkpoint
2. zero or more checkpoints that are not stable
3. current state

proof of correctness for a checkpoint, is generated as follows:
>
1. when a replica i produces a checkpoint, it multicasts a message (CHECKPOINT, n, d, i) to the other replicas, 
>> n is the sequence number of the last request whose execution is reflected in the state, <br>
>> d is the digest of the state.
2. each replica collects checkpoint messages in its log until it has 2f+1 of them for sequence number n with the same digest d signed by different replicas.   These 2f+1 messages are the proof of correctness for the checkpoint.

A checkpoint with a proof becomes stable and the replica discards all pre-prepare, prepare and comit messages with sequence number less than or equal to n from its log; it also discards all earlier checkpoints and checkpiont messages.

The checkpoint is used to advance the low and high water marks.  The low-water mark h is equal to the sequence number of the last stable checkpoint.  The high water mark H = h + k.


### View Changes

View changes are triggered by timeouts that prevent backups from waiting indefinitely for requests to execute.  A backup is waiting for a request if it received a valid request and has not executed it.  A backup starts a timer when it receives a request and the timer is not already running.  It stops the timer when it is no long waiting to execute the request, and restarts it if at that point it is waiting to execute some other request.


***initiate view change***
>
1. backup starts a view change to move the system to view v+1
2. it stops accepting messages (other than checkpoint, view-change, and new-view message)
3. multicast a (VIEW-CHANGE, v+1, n, C, P, i)
>> n: the sequence number of the last stable checkpoint s known to i <br>
>> C: a set of 2f+1 valid checkpoint message proving the correctness of checkpoint s <br>
>> P: a set of Pm for each request m that prepared at i with a sequence number higher than n <br>
>> each Pm contains a valid pre-prepare message and 2f matching
4. primary p of view v+1 recieves 2f valid view-change messages for view v+1 from other replicas.
5. primary multicasts (NEW-VIEW, v+1, V, O) message to all other replica
>> The primary determines the sequence number min-s of the last stable checkpoint in V and the highest sequence number max-s in a prepare message in V <br>
>> primary creates a new pre-prepare message for view v+1 for each sequence number n between min-s and max-s. <br>
>> V: a set of containing the valid view-change message received by the primary plus the view-change message for v+1 the primary sent <br>
>> O: a set of pre-prepared messages, from min-s to max-s
6. primary appends the message in O to its log.
>> If min-s is greater than the sequence number of its latest stable checkpoint, the primary also inserts the proof of stability for the checkpoint with sequence number min-s in its log.
7. enters view v+1
8. at this point, it is able to accept message for view v+1.


***accept view change***

A backup accepts a new-view message for view v+1, iff
>
1. the msg is signed correctly
2. the view-change messages it contains are valid for view v+1
3. the set O is correct
4. add the new information to its log as described for the primary
5. multicast a prepare for each message in O to all the other replicas, add these prepares to its log
6. enter view v+1

replicas redo the protocol for the messages between min-s and max-s, but they avoid re-executing client requests.


