


P : the set of proposers

A : the set of acceptors

B : the set of proposal numbers, also called ballot numbers, which is any set that can be strictly totally ordered.

S : a set of slots used to index the sequence of decisions.

Q : a subset of the acceptors, is used as a quorum system

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

definition Phase1a :: "b \<in> B"
                where
                "phase1a b = 
                        <if>   \<not> (<exists> m. m \<in> msgs \<and> m.type \<eq> '1a' \<and> m.bal \<eq> b)
                        <then> send(['type': '1a', 'bal': b])
                               \<and> unchanged  (maxVBal, maxBal, maxVal)


definition Phase1b :: "a \<in> A"
                      where
                      "phase1b a = 
                        \<if> \<exists> m \<in> msgs 
                              \<and> m.type \<eq> "1a"
                              \<and> m.bal \<gr> maxVal[a]
                        \<then> send(['type': '1b', 'bal': m.bal, 'maxVBal': maxVBal[a], 'maxVal': maxVal[a], 'acc': a])
                                \<and> maxBal' = maxBal \<if> a \<ne> m.bal
                                \<and> unchanged (maxVBal, maxVal)



definition Phase2a :: "b \<in> B"
                      where
                      "phase2a a = 
                      \<if> \<not> \<exists> m. (m \<in> msgs \<and> m.type \<eq> '2a' \<and> m.bal \<eq> b)
                            \<and> \<exists> v. v \<in> V 
                            \<and> \<exists> q \<in> Q
                            \<and> \<exists> S \<subset> { \<forall> m. m.type \<eq> '1b' \<and> m.bal \<eq> b }
                            \<and> \<forall> a. a \<in> Q \<and> \<exists> m \<in> S \<and> m.acc \<eq> a
                            \<and> (\<forall> m. m \<in> S \<and> m.maxVBal <eq> -1
                                    \<or> \<exists> c. c \<in> {0 ... (b-1)}
                                          \<and> (\<forall> m. m \<in> S \<and> (m.maxVBal \<le> c \<or> (m.maxVal \<eq> c \<and> m.maxVal <eq> v))))
                      \<then> send(['type': '2a', 'bal': b, 'val': v])
                              \<and> unchanged(maxBal, maxVBal, maxVal)



definition Phase2b :: "a \<in> A"
                      where
                      \<if> \<exists> m. m \<in> msgs:
                            \<and> m.type \<eq> '2a'
                            \<and> m.bal \<ge> maxBal[a]
                      \<then> send(['type': '2b', 'bal': m.bal, 'val': m.val, 'acc': a])
                              \<and> maxBal' = maxBal \<if> a \<ne> m.bal
                              \<and> maxVBal' = maxBal \<if> a \<ne> m.bal
                              \<and> maxVal = maxVal \<if> a \<ne> m.val





