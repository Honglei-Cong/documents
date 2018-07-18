## HyParView

only rely on partial views instead of complete membership information.

1.	gossip strategy that is based on the use of a reliable transport protocol to gossip between peers.
2.	each node maintains a small symmetric active view the size of the fanout+1.  Broadcast is performed by flooding the graph defined by the active views.
3.	TCP is also used as a failure detector, and since all members of the active view are tested at each gossip step
4.	Each node maintains a passive view of backup nodes that can be promoted to the active view, when one of the nodes in the active view failed.
5.	membership protocol is in charge of maintaining the passive view and selecting which members of the passive view should be promoted to the active view.


### Motivation

* The fanout of a gossip protocol is constrained by the target reachability level and the desired fault tolerance of the protocol
* Gossip would strongly benefit from fast healing properties.  High failure rates may have a strong impact on the quality of partial views.

Fanout: tradeoff of reliability and message duplication.

Fast Healing:


##### HyParView

```
c = 1, k = 6, 
active view: log(n) + c
passive view: k * (log(n) + c)
```

* Links in the overlay are symmetric
* When a node receives a message for the first time, it broadcast the message to all nodes of its active view.
* reactive strategy is used to maintain the active view.  The reader should notice that each node tests its entire active view every time it forwards a message.  Therefore, the entire broadcast overlay is implicitly tested at every broadcast, which allows a very fast failure detection.
* Passive view is not used for message dissemination.  The goal of the passive view is to maintain a list of nodes that can be used to replace failed members of the active view.
* Periodically, each node performs a shuffle operation with one of its neighbors in order to update its passive view.

Identifiers are exchanged in a shuffle operation are not only from the passive view, a node also sends its own identifier and some nodes collected from its active view to its neighbors.  This increases the probability of having nodes that are active in the passive views and ensures that failed nodes are eventually expunged from all passive views.



### Join

* Join
* ForwardJoin
* Disconnect

message TTL:

* Active Random Walk Length
* Passive Random Walk Length


### Active View Management

1.	reactive strategy
2.	NEIGHBOR request


### Passive View Management

1.	SHUFFLE request = (p, K_a nodes from active view, K_p nodes from passive view)
2.	SHUFFLEREPLAY response


