
## Introduction

### features of Raft

>
* strong leader: log entries only flow from the leader to other servers.  This simplifies the management of the replicated log and make it easier to understandy.
* leader election: uses randomized timers to elect leaders.
* membership changes: uses a new **joint consensus** approach where the majorities of two different configurations overlap during transitions.  This allows the cluster to continue operating normally during configuration changes.

### replicated state machine

consensus algorithms for practical systems typically have the following properties:
>
* safety under all non-byzatine conditions, including network delays, partitions, and packet loss, duplications, and reordering.
* available as long as any majority of the servers are operational and can communicate with each other and with clients.
* Servers are assumed to fail by stopping: they may later recover from state on stable storage and rejoin the cluster.
* Do not depend on timing to ensure the consistency of the logs: faulty clocks and extreme message delays can, at worst, cause availability problems.
* A command can complete as long as majority of the cluster has responded to a single round of remote procedure calls; a moinority of slow servers need not impact overall system performance.

### Algorithm

Raft implements consensus by first electing a distinguished leader, then giving the leader complete responsibility for managing the replicated log.  The leader accepts log entries from clients, replicates them on other servers, and tells servers when it it safe to apply log entries to their state machines.

Raft decomposes the consensus problem into three relative independent subproblems:
>
* Leader election: new leader must be chosen when an existing leader fails.
* Log replication: the leader must accept log entries from client, and replicate them across the cluster, forcing the other logs to agree with its own.
* Safety: If any server applied a particular log entry to its state machine, then no other server may apply a different command for the same log index.

## Raft basics

At any time, server in one of three states: Leader, Follower, or Candidate.

### Term
Raft divides time into **terms** of arbitrary length.  Terms are numbered with consecutive integers.  Each term begins with an election, in which one or more candidates attempt to become leader.

Term act as a logical clock in Raft, and they allow servers to detect obsolete information such as stale leaders.  Each server stores a **current term** number, which increases monotonically over time.

Current terms are exchanged whenever server communicate; if one server's current term is smaller than the other's, then it updates its current term to the larger value.  If a candidate or leader discovers that its term is out of date, it immediately reverts to follower state.  If a server receives a request with a stale term number, it rejects the request.

Basic Raft only have two types of RPCs:

>
* RequestVote 

>>
>> 1. term: candidate's term
>> 2. candidateId: candidate requesting vote
>> 3. lastLogIndex: index of candidate's last log entry
>> 4. lastLogTerm: term of candidate's last log entry

>
* AppendEntries

>>
>> 1. term: leader's term
>> 2. leaderId: so followers can redirect clients
>> 3. prevLogIndex: index of log entry immediately preceding new ones
>> 4. entries[]: log entries to store (empty for heartbeat)
>> 5. leaderCommit: leader's commitIndex


## Leader election

Raft uses a heartbeat mechanism to trigger leader election.  When servers start up, they begin as followers.  A server remains in follower state as long as it receives valid RPCs from a leader or candidate.  Leaders send periodic heartbeats (AppendEntries RPC that carry no log entries) to all followers in order to maintain authority.  If a follower receives no commnunication over a period of time (election timeout), then it assumes there is no viable leader and begins an election to choose a new leader.

procedures:
>
1. follower increments its current term, and transitions to candidate state.
2. follower vote for itself, and issue RequestVote RPC in parallel to each of the other servers in the cluster.
3. A candidate continues in the state until:
>> a. it wins the election
>> b. another server establishes itself as leader
>> c. a period of time goes by with no winner.
4. candidate wins the election when it receives votes from a majority of servers for the same term. (each server will vote for at most one candidate in a given term, on a First-come-first-service mode).
5. candidate becomes the leader.
6. leader sends heartbeat messages to all of the other servers to establish its authority and prevent new elections.


Exception (1)
>
4.1 which waiting for votes, candidate may recieve an AppendEntries RPC from another server claiming to be leader
4.1.a  If the leader's term in its RPC is at least as large as the candidate's current term, then the candidate recognizes the leader as legitimate and returns to follower state
4.1.b  If the term in RPC is smaller than the candidate's current term, then the candidate rejects the RPC and continues in candidate state.
4.1.c  If multiple candidate timeout and start a new election by incremeting its term and initiating another round of RequestVote RPCs.  Raft uses randomized election timeout to ensure that split votes are rare and that they are resolved quickly.


## Log Replication

Leader applies the entry to its state machine and returns the result of that execution to the client.  If followers crash or run slowly, or if network packets are lost, the leader retries AppendEntries RPC indefinitely (even after it has responded to the client) until all followers eventually store all log entries.

Each log entry also has an integer index identifying its position in the log.  The leader decides when it is safe to apply a log entry to the state machines: such an entry is called **committed**.  The leader keeps track of the highest index it knows to be committed, and it includes that index in future AppendEntries RPCs (including heartbeats) so that the other servers eventually find out.  Once a follower learns that a log entry is committed, it applies the entry to its local state machine (in log order).

A follower may be missing entries that are present on the leader; it may have extra entries that are not present on the leader, or both.  Missing and extraneous entries in log may span multiple terms.  In Raft, the leader handles inconsistencies by forcing the follower's log to duplicate its own.  

To bring a follower's log into consistency with its own, the leader must find the latest log entry where the two logs agree, delete any entries in the follower's log after that point, and send the follower all of the leader's entries after that point.  All of these actions happen in response to the consistency check performed by AppendEntries RPCs.  The leader maintains a **nextIndex** for each follower, which is the index of the next log entry the leader will send to that follower.

how to AppendEntries:
>
1. Leader first comes to power, it initialize all nextIndex values to the index just after the last one in its log.
2. If follower's log is inconsistent with the leader's, the AppendEntries consistency check will fail in next AppendEntries
3. after a rejection, leader decrement nextIndex and retries the AppendEntries
4. Eventually, nextIndex will reach a point where the leader and follower logs match.


## Safety

### Election Restriction

Raft uses a simpler approach where it guarantees that all the committed entries from previous terms are present on each new leader from the moment of its election, without the need to transfer those entries to the leader.

Raft uses the voting process to prevent a candidate from winning an election unless its log contains all committed entries.  A candidate must contact a majority of the cluster in order to be elected, which means that every committed entry must be present in at least one of those servers.  

The RequestVote RPC implements its restrictions:  **The RPC includes information about the candidate's log, and the voter denies its vote if its own log is more up-to-date than that of the candidate.**

### Committing entries from previous terms

A leader knows that an entry from its current term is committed once that entry is stored on a majority of the servers.  If a leader crashes before committing an entry, future leaders will attempt to finish replicating the entry.  

However, a leader cannot immediately conclude that an entry from a previous term is committed once it is stored on a majority of servers.  Raft never commits logs entries from previous terms by counting replicas.  Only log entries from the leader's current term are committed by counting replicas; once an entry from the current term has been committed indirectly in this way, then all prior entries are committed indirectly because of the Log Matching Property.


## Cluster Membership changes

In order to ensure safety, configuration changes must use a two-phase approach.  There are a variety of ways to implement the two phases.  For example, some systems use the first phase to disable  the old configuration so it can not process client requests; then the second phase enables the new configuration.  In Raft the cluster first switches to a transitional configuration we call **joint consensus**; once the joint consensus has been committed, the system then transitions to the new configuration.
>
1. Log entries are replicated to all servers in both configurations.
2. Any server from earlier configuration may server as leader.
3. Agreement (for elections and entry commitment) requires separate majorities from both the old and new configuraitons.

The joint consensus allows individual servers to transition between configurations as different times without compromising safety. Furthermore, joint consensus allows the cluster to continue servicing client requests throughout the configuration change.

Process:
>
1. leader recieves a request to change the configuration from C-old to C-new
2. leader stores the configuraiton for joint consensus as a log entry
3. leader replicates that entry using AppendEntries
4. Servers adds the new configuration entry to its log
5. Servers use that ocnfiguraiton for all future decisions

This means that the leader will use the rules of C-old.new to determine when the log entry for C-old.new is committed.

Once C-old.new has been committed, either C-old nor C-new can make decisions without approval of the other, and the Leader Completeness Property ensures that the only servers with the C-old.new log entry can be elected as leader.  It's now safe for the leader to create a log entry describing C-new and replicate it to the cluster, and this configuration will take effect on each server as soon as it is seen.  When the new configuration has been committed under the rule of C-new, the old configuration is irrelevant and servers not in the new configuration can be shut down.

Three more issues:
>
1. new servers may not initially store any log entries.  
>> It will take quite a while for them to catch up, during which time it might not be possible to commit new log entries.  In order to avoid availability gaps, Raft introduces an additional phase before the configuration change, in which the new servers join the cluster as non-voting members (the leader replicates log entries to them, but they are not considered for majorities).
2. cluster leader may not be part of new configuration.
>> In this case, the leader step down (returns to follower state) once it has committed the C-new log entry.
3. removed servers (those not in C-new) can disrupt the cluster.  
>> These servers will not receive heartbeats, so they will time out and start new elections.  To prevent this problem, servers disregard RequestVote RPCs when they believe a current leader exists.  Specially, if a server receives a RequestVote RPC within the minimum election timeout of hearing from a cluster leader, it does not update its term or grant its vote.

## Log Compaction

Snapshot is the simplest approach to compaction.  In snapshot, the entire current system state is written to a snapshot on stable storage, then the entire log up to that point is discarded.  Snapshot is used in Chubby and ZooKeeper.

InstallSnapshot RPC:
>
Arguments:
>> * term: leader's term
>> * leaderId: so follower can redirect clients
>> * lastIncludedIndex: the snapshot replaces all entries up through and including this index
>> * lastIncludedTerm: term of lastIncludedIndex
>> * offset: byte offset where chunk is positioned in the snapshot file
>> * data[]: raw bytes of the snapshot chunk, starting at offset
>> * done: true if this is the last chunk
>
Procedures:
>> 1. reply immediately if term < currentTerm
>> 2. create new snapshot file if first chunk (offset is 0)
>> 3. write data into snapshot file at given offset
>> 4. reply and wait for more data chunks if done is false
>> 5. save snapshot file, discard any existing or partial snapshot with a smaller index
>> 6. if existing log entry has same index and term as snapshot's last included entry, retain log entries following it and reply
>> 7. discard the entire log
>> 8. reset state machine using snapshot contents (and load snapshot's cluster configuration).

## Client Interaction

Clients of Raft send all of their requests to the leader.  When client starts, it connects to a randomly chosen server.  If the chosen server is not leader, that server will reject the client's request and supply information about the most recent leader it has heard from.  

If the leader crashes, client request will be timed out; client then try again with randomly-chosen servers.

To ensure linearizable semantics, client assign unique serial numbers to every command.  Then, the state machine tracks the latest serial number processed for each client, along with the associated response.  If it receives a command whose serial number has already been executed, it responds immediately without re-executing the request.

To prevent returing stale data for read-only operations, linearizable reads must not return stale data, and Raft needs two extra precautions to guarantee this without using the log:
>
1. A leader must have the latest information on which entries are committed.
2. At the start of leader's term, it may not know which those are.  To find out, it needs to commit an entry from its term.  Raft handles this by having each leader commit a blank no-op entry into the log at the start of its term.
3. Leader must check whether it has been deposed before processing a read-only request.  Raft handles this by having the leader exchange heartbeat messages with a majority of the cluster before responding to read-only requests.  Alternatively, the leader could rely on the heartbeat mechanism to provide a form of lease, but this would rely on timing for safety (it assumes bounded clock skew).

