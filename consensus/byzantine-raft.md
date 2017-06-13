
## Features of BFT Raft

>
1. Message Signatures
>> BFT Raft uses digital signatures extensively to authenticate messages and verify their integrity.  This prevents a Byzantine leader from modifying the message contents or forging messages. Client public keys are kept separate from replica public keys to enforce that only clients can send new valid commands, and only replicas can send valid Raft RPCs.
2. Client intervention
>> BFT Raft allows clients to interrupt the current leadership if it fails to make progress.  This allows BFT Raft to prevent Byzantine leaders from starving the system.
3. Incremental hashing
>> Each replica in BFT Raft computes a cryptographic hash every time it appends a new entry to its log.  A node can sign its last hash to prove that it has replicated the entirety of a log, and other server can verity this quickly using the signature and the hash.
4. Election verification
>> Once a node has become leader, its first AppendEntries RPC to each other node will contain a quorum of RequestVoteResponse RPCs that it received.  Nodes first verify that the leader has actually won an election by counting and validating each of these RequestVoteResponse.  Subsequent AppendEntries RPCs in the same term to the same node need not includes these votes, but a node can ask for them in its AppendEntriesResponse if necessary.
5. Commit verification
>> In order to safely commit entries, each AppendEntriesResponse RPC is broadcast to each other node, rather than just to the leader.  Further, each node decides for itself when to increment its commit index, rather than the leader.  Once a node has received a quorum of matching AppendEntriesResponse RPCs from other nodes at a particular index, it considers that to be the new commit index, and discards any stored AppendEntriesResponse RPCs for previous indices.
>> AppendEntriesResponse RPCs becomes closer to the broadcast PREPARE message from PBFT.
6. Lazy Voters
>> A node does not grant a vote to a candidate unless it believes the current leader is faulty.  A nodes comes to believe that its leader is faulty if it does not receive an AppendEntries RPC within its own election timeout, or it recevies an **UpdateLeader RPC from a client** for that leader.

## BFT Raft Algorithm

safety grarantees as Raft:
>
* Election safety
* Leader Append-Only
* Log Matching
* Leader Completeness
* State Machine Safety

### BFT Raft Basics

A BFT Raft cluster that tolerates f Byzantine failures must contain at least (n >= 3f + 1), where (n-f) nodes form a quorum.  BFT Raft is configured so that nodes and clients have the public keys of each other node and client ahead of time.  Nodes and clients always sign before sending messages and reject messages that do not include a valid signature.

Raft RPCs are idempotent.  

For clients commands, we use a monotonically increasing per-client indentifier to detect duplicated messages.

Compared with Raft, BFT Raft nodes only update their term number in one of three situations:
>
1. On receiving an AppendEntries that contains a quorum of votes for the sender in a new term
2. when responding to a RequestVote for a higher term
3. when becoming a candidate.

The consensus algorithm requires four type of RPCs:
>
1. AppendEntries
>> initiated by leaders to replicate log, provide a heartbeat, and communicate a successful election
2. RequestVote
>> initiated by candidates during election
3. SendRequest
>> initiated by clients to request the cluster to execute a command to its replicated state machine
4. UpdateLeader
>> initiated by clients to request change of leadership

### Incremental Hashing

### Log Replication

As in PBFT, clients need to wait for (f+1) matching replies to each request before exposing that result to application logic.  This ensures that at least one honest node decided that a particular result was safely replicated to a quorum and should be externalized.

Replication:
> 
1. As in Raft, the node will check that it has a matching log prefix, but using the incremental hash, rather than the term of the previous entry.
2. It then check the authenticity of the of the new entries for itself.
3. If all valid, the node will append new entries to its log, and compute the incremental hash at each new index.
4. It will then broadcast its AppendEntriesResponse to each other node, which contains the incremental hash at the last new index.

Commit:
>
1. When a node receives An AppendEntriesResponse, it will save it if it is for an index higher than the nodes current commit index.
2. Once a node receives a quorum of matching AppendEntriesResponse for a particular index, it is afe to commit everything up to that log entry.
3. The node then apply the new committed log entries to its state machine
4. and send results to the client directly.
5. nodes also store the results of the committed entries, so they can be retransmitted if need.

for Leader, if it receives an AppendEntriesResponse, in addition, leader will check if the AppendEntries was successful.  The node could have responded that :
>
1. it needed proof of that leader's successful election, 
>> Leader will send a new AppendEntries RPC containing the votes it got during the election.
2. or that it did not have the previous entry to the new entries.
>> Leader decrement the nextIndex for that node and retry, just as in Raft.

### Leader Election

In addition to the spontaneous follower-triggered elections, BFT Raft also allow client intervention: when a client observes no progress with a leader for a period of time, called **progress timeout**, it broadcast UpdateLeader RPCs to all nodes, telling them to ignore future heartbeats from what the client believes to be the current leader in the current term.

When a node receives a RequestVote RPC with a valid signature, it grants a vote only if all five conditions are true:
>
1. the node has not handled a heartbeat from its current leader within its own timeout
2. the new term is between its (current term + 1) and (current term + H)
3. the request sender is an eligible candidate
4. the node has not voted for another leader for the proposed term
5. the candidate shares a log prefix with the node that contains all committed entries.








