
## Asynchronous Broadcast Protocols


**secure causal atomic broadcast**
This is a robust atomic broadcast protocol that tolerates a Byzantine adversary and also provides secrecy for messages up to the moment at which they are guaranteed to be delivered.
Thus, client requests to a trusted service using this broadcast remain confidential until they are answered by the service and the service processes the requests in a causal order.

This is done by combining:
* atomic broadcast protocol
* threshold decryption

#### Layered architecture

* Secure Causal Atomic Broadcast 
* Atomic Broadcast
* Multi-valued Byzantine Agreement
* Broadcat Primitives + Byzantine Agreement


### Validated Byzantine Agreement (VBA)

generalizes the primitives of *agreement on a core set*

(ID, in, v-propose, v), where v in {0, 1}*
(ID, out, v-decide, v)


The basic idea is that every party proposes its value as a candidate value for the final result.  One party whose proposal satisfies the validation predicate is then selected in a sequence of binary Byzantine agreement protocols and this value becomes the final decision value.

1. Echoing the proposal
   1. Each party *c-broadcast* the value that is proposes to all other parties using verifiable authenticated consistent broadcast.
   2. Then `Pi` waits until it has received `n-f` proposals satisfying Q_{ID} before entering the agreement loop
2. Agreement Loop
   One party is chosen after another, according to a fixed permutation.  Each party `Pi` carries out the following steps for `Pa`:
   1. send a *v-vote* to all parties containing 1 if `Pi` has received `Pa`'s proposal and 0 otherwise
   2. wait for `n-f` *v-vote* message, but do not count votes indicating 1 unless a valid proposal from `Pa` has been received (either directly or included in the *v-vote* message)
   3. run a binary validated byzantine agreement biased towards to 1 to determine whether `Pa` has properly broadcast a valid proposal.
3. Delivering the chosen proposal


### A Constant-round Protocol for Multi-valued Agreement

choose the order randomly during the protocol after making sure that enough parties are already committed to their votes on the candidates.

This is achieved in two steps:
1. each party must commit to the votes that it will cast by broadcasting the identities of the `n-f` parties from which it has received valid *v-echo* messages.
   1. Honest parties will later only accept *v-vote* messages that are consistent with these commitments
2. determine the permutation using a threshold coin-tossing scheme that outputs a pesudo-random value, after enough votes are committed.


### Atomic Broadcast

uses a digital signature scheme *`S`*

1. Each party maintains a FIFO queue of not yet *a-delivered* payload messages.
   1. messages recieved to *a-broadcast* are appended to this queue whenever they are received

The protocol proceeds in a asynchronous global rounds, where each round *r* consists of the following steps:

1. send the first payload message *`w`* in the current queue to all parties, accompanied a digital signature *`s`* in an *a-queue* message
2. collect the *a-queue* message of `n-f` parties and store them in a vector *`W`*, and propose *`W`* for validated byzantine agreement.
3. Perform multi-valued byzantine agreement with validation of a vector of tuples `W = [(w1, s1), ... (wn, sn)]` 
4. after deciding on a vector *`V`* of messages, deliver the union of all payload messages in *`V`* according to a deterministic order; proceed to next round.

