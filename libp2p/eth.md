## Eth Project

https://github.com/libp2p/libp2p/issues/33


Furthermore, with Ethereum 2 research beginning as a relatively green-field project, they would like to identify the most stable/performant pieces of libp2p to maximize ease of development vs future pain.

* Swappable transports
* Zero round-trip handshakes (a la DTLS/QUIC)
* Stronger encrypted transports
* Attack mitigation (DDOS, timing attacks, etc. Obviously a big one!)
  * Security audits planned/executed?
  * Node identity and authentication (perhaps verifiable from sources outside of peer network)
  * Peer selection strategies (is random-walk appropriate)
o	In degree vs out degree measurement for connection rebalancing
* Connection management (likely relevant to ongoing peer swarm refactoring)
  * Pluggable ConnManager?
* Gossipsub
  * Performance characteristics
  * What are its IPFS/PL use cases? Will this implementation stick around and get the love it needs to become a core part of either Ethereum or PL infrastructure.
* Potential IPFS/Ethereum interactions
* Python/other language support
  * libp2p daemon is likely highest priority

  
### Sharding-specific Interests
* DHT persistence, speed of bootstrapping
  * When rejoining the network, are your peers good/connected?
  * DHT routing table persistence not yet landed
* How might we map shards to libp2p?
  * Separate peer network per shard
  * Single channel that dispatches clients to separate peer networks
  * Multiple channels overlayed over a single shared peer network
* How can we optimize balance of flood/gossipsub propagation vs active connections?
* Support for authentication in flood/gossipsub?
* Fast shard jumping
  * Can we (perhaps at a slower interval) keep a small cache of DHT routing table for other shards?
  * Can we keep an open connection to some small number of peers in shards we're not actively participating in to speed this switch
  * QUIC transport clearly a huge boon here


I'll edit this initial post as I get more clarity on requirements and requests!

