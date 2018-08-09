## cOS


* off-chain decentralized application operating system
* a combination of SDK and runtime system

cOS models a system of off-chain applications as a directed-acyclic-graph of conditionally dependent states, where the edges represent the dependencies among them.


#### Off-chain application development framework
In general, we categorise decentralized applications into two classes:
 
* simple pay-per-use application
    * user keeps receiving microsevices (e.g. data relay) from a real-world entity and streams payments through the payment network
* multi-party applications
    * metaprogramming
    * annotation processing
    * dependency injection
    * compiler process the application code, extracts the declared off-chain objects, and generate the conditional dependency graph.
    * compiler detects invalid or unfulfillable dependency information and generates human-readable errors to assist the developer in debugging.

#### Off-chain application runtime

cOS servers as the interface between cApps and the Celer Network transport layer.  It supports cApps in terms of both network communication and local off-chain state management.

**On the network front**
the runtime handles multi-party communication during the lifecycle of a cApp.  It also provides a set of primitives for secure multi-party computation capable of supporting complex use cases such as gaming.  

In case of counter-party failure, the runtime relays disputes to the on-chain state.  In the case of client going offline, the runtime handles availability offloading to the State Guardian Network.  When the client comes back online, the runtime synchronies the local stats with the Start Guardian Network.

**At its core**
cOS runtime bundles a native virtual machine for running smart contract.

* VM progresses through the same bytecode as if they were executed on-chain
* VM needs to update and store the states locally instead of on the blockchain
* VM can shut down unexpectedly at any time due to software bug.  To avoid corruption of local states, we need to implement a robust logging, checkpointing and committing protocol.
* gas metering can be omitted, because the execution happens locally and it does not make sense to change gas fees.


