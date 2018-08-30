


P : the set of proposers

A : the set of acceptors

B : the set of proposal numbers, also called ballot numbers, which is any set that can be strictly totally ordered.

S : a set of slots used to index the sequence of decisions.

V : the set of possible proposed values

phi : a predicate that given a value v evaluates to true iff v was chosen by the algorithm


### Phase1

a : a proposer selects a proposal number n and sends a prepare request with number n to a majority of acceptors.

b : If an acceptor receives a prepare request with number n greater than that of any prepare request to which it has already responded, then it responds to the request with a promise not the accept any more proposals numbered less than n and with the highest-numbered proposal (if any) that it has accepted.


### Phase2

a : If the proposer receives a response to its prepare requests (numbered n) from a majority of acceptors, then it sens an accept request to each of those acceptors for a proposal numbered n with a value v, where v is the value of the highest-numbered proposal among the responses, or is any value if the responses reported no proposals.


b : If an acceptor receives an accept request for a proposal numbered n, it accepts the proposal unless it has already responsed to a prepare request having a number greater than n.




### Specifications


msgs - the set of messages that have been sent.

accVoted - per acceptor, a set of triples in B x S x V, capturing a mapping in S -> B x V, that the accepor has voted for.  This contrasts two numbers per acceptor, in two variables, maxVBal and maxVal, in the specification of Basic Paxos.

accMaxBal - per acceptor, the highest ballot number seen by the acceptor.  This is named maxBal in the specification of Basic Paxos.

proBallot - per proposer, the ballot number of the current ballot being run by the proposer.  This is added to support preemption.



#### Basic Paxos

Phase1a



