## Requirement

a draft for next p2p server design.


1.	support tcp/udp
2.	support multi-network
3.	support broadcast/unicast
4.	gossip based DHT(epidemic broadcast tree + GoCast)
5.	support ip-address filtering policies
6.	support membership management (HyParView)
    1. dynamic membership
7.	peer auto-discovery
8.	all peers are authenticated
9.	all msgs are authenticated
10.	based on protobuf (for interoperations)
    1. backward compatible
11.	support connection-based/peer-based metrics
12.	support connection-based-throttling
13.	support connection-accept-limitation
14.	each peer only accepts limited(configured) in-degree/out-degree
15.	persistent peer information, to support quick recovery/bootstrap
16.	support peer liveness detection (ping)
17.	10000+ node network, with 10+ seeds
18.	network failure resistance 
    1. <10% node failures, bring <5% performance penalty
    2. provable resistant on >50% node failures
19.	reactor-based(event-driven) design
20.	push based (combined eager-push and lazy-push)
21.	support network simulation
22.	support NAT (optional)


### Reference

* Epidemic Broadcast Tree
* GoCast: Gossip-Enhanced Overlay Multicast for Fast and Dependable Group Communication




