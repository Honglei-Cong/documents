# Parity Substate

https://github.com/paritytech/substrate
https://slides.com/paritytech/paritysubstrate#/
https://slides.com/paritytech/substrate_web3summit#/

A framework to build new blockchains.

Polkadot is built with Substate and projects built with Substate can run natively on Polkadot.

### What is Substate

An application framework to build distributed or decentralised system such as cryptocurrencies or a message bus.


### What Substate provides

To build a new project with Substate, all you have to do is implement a very small number of hooks in your code.

Substate provides:

* Consensus, finality and block voting logic.
* Networking, so peer discovery, replication etc.
* An efficient, deterministic, sandboxed WebAssembly runtime.
* The ability to seamlessly run a node in the browser that can communicate with any desktop and cloud node.
* A cross-platform database/file storage abstraction
* Seamless updates.
* The ability to immediately start running your project on Polkadot.


### How to use Substate

You need provide:

* your state machine, which includes things like transactions.

Use Substate to create a decentralised Erlang-style actor-model concurrent system with a set of trusted authorities to verify the correct behavior of the network.

To get a full blockchain up and running, you'd need to implement:

* A function that creates a new pending block based on the previous block's header.  The header includes:
    * block height
    * a 'crypto commitment' to the block's state
    * a crypto commit to all the extrinsics in the body
    * a hash of the block's parent
    * some extra arbitrary data
* A function that adds an extrinsic to a pending block.  This should also update the chain's state (e.g. account balances)
* A function that takes a pending block and generates a finished block from it.
* A function that executes the existing block.
    
