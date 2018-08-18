## cRoute

to build a network consisting of state channels, while state transitions should be relayed in a trust-free manner.

* how fast and how many transactions can flow on a given network

**shortest-path-routing**

* stateful link model
* does not account for channel balancing, and thus each link may quickly run out of its capacity, which further leads to frequent changes in network ontology.

**Distributed Balanced Routing (DBR)**

* inspired by the BackPressure routing algorithm
* the routing decision is guided by the current network’s congestion gradients
* provably optimal throughput
* transparent channel balancing
* fully decentralised
    * each node only needs to talk to its neighbours in the state channel network topology.
* failure resilience
    * quickly detect and adapt to unresponsive nodes, supporting the maximum possible throughtput over the remaining available nodes.
* privacy preserving

**Debt Queue**

* each node needs to maintain a ‘debt queue’ for payments destined to each node k, whose queue length corresponds to the amount of tokens that should be relayed by node i to the next hop but have not been relayed yet at the beginning of slot t.
* Intuitively, the length of the debt queue is an indicator of congestion over each link.

**Channel Imbalance**

* channel imbalance is the difference between the total amount of tokens received by node i from node j and the total amount of tokens sent from i to j over their payment channel up to the beginning of slot t.

**Congestion-plus-imbalance (CPI) weight**

* CPI weight for link i->j  and destination k, is the sum of differential backlog Qi - Qj for payments destined to node k between node i and node j and the channel imbalance between node i and node j.


during each slot t, each node i receives new payment requests from outside the network, where the total amount of tokens that should be delivered to node k is a_i_k(t) > 0.  Also denote by u_ij_k(t) the amount of tokens (required to delivered to node k) sent over link i->j in slot t, which is referred to as a routing variable.


