## Epidemic Broadcast Tree

#### Tree Construction
selecting which links of the random overlay network will be used for forward message payload using an eager push stategy.

#### Tree Repair
Repair the tree when failure occurs.
It should be able to detect and heal partitions of the tree.

#### Peer Sampling Service
expose a GetPeers() primitive in its interface, which is used by the protocol to get information about neighbours to whom it should send messages.

#### Exported Neighbor Interface
NeighborUp()/NeighborDown() interfaces are exposed from protocol, to notify protocol whenever a change happens on the partial view maintained by the peer sampling service. These primitives are used to support quick healing of the broadcast tree.

#### Process

1.	in order to broadcast a message, each node gossips with F nodes provided by a peer sampling service
2.	Eager push (EagerPushPeers) is used just for a subset of F nodes; lazy push is used for the remaining nodes.
3.	Eager push (LazyPushPeers) links are selected in such a way that their closure effectively builds a broadcast tree.
4.	The set of peers are not changed at each gossip round until failures are detected.


#### Peer Sampling Service Requirements

1.	Connectivity
2.	Scalable
3.	Reactive Membership
4.	Symmetric partial view


#### Gossip and Tree Construction

1.	Initialize EagerPush peers: contains F random peers initially, obtained through peer sampling service.
2.	Construct the spanning tree, by moving neighbors from EagerPushPeers to LazyPushPeers.
3.	When a node receives a message for the first time, it includes the sender in the set of EagerPushPeers
4.	When a duplicated is received, its sender is moved to LazyPushPeers.  Furthermore, a PRUNE message is sent to that sender such that, in response, it also moves the link to LazyPushPeers.
5.	Assuming a stable network, it will tend to generate a spanning tree that minimizes the message latency (any provement)
6.	Lazy push is implemented by sending IHAVE messages, that only contain the broadcast ID, to all LazyPushPeers.
7.	IHAVE do not need to be sent immediately.  A scheduling policy is used to piggyback multiple IHAVE announcements in a single control message.


#### Tree Repair

1.	When received IHAVE, it simply mark the corresponding message as missing.  It then start a timer, and waits for the missing message to be received via eager push before the timer expires.
2.	When the timer expires, that node selects the first IHAVE announcement it has received for the missing message, and send GRAFT message.
a)	GRAFT message triggers the transmission of the missing message payload.
b)	add the corresponding link to the broadcast tree.
3.	Start another timer after GRAFT sent, to ensure that the message will be requested to another neighbor if it is not received meanwhile.



### Dynamic Membership

#### NeighborDown

When a neighbor is detected to leave the overlay, it is simply removed from the membership.  The record of IHAVE messages sent from failed members is deleted from the missing history.

#### NeighborUp

When a new member is detected, it is simply added to the set of EagerPushPeers, i.e., it is considered as a candidate to become part of the tree.



## Messages

#### GOSSIP
        
        * eager push message

#### IHAVE

        * Lazy pull message

#### PRUNE

        * to remove one active link

#### GRAFT

        * triggers the transmission of the missing message payload
        * adds the corresponding link to the broadcast tree, healing it

#### NEIGHBOR-UP

        * peer sampling service notify DHT that new neighbor discovered.

#### NEIGHBOR-DOWN

        * peer sampling service notify DHT that neighbor disconnected.




