
## Plasma

five key components:

* An incentive layer
    * for persistently computing contracts in an economically efficient manner
* structure for arranging child chains in a tree format
    * to maximize low-cost efficiency and net-settlement of transactions
* a MapReduce computing framework
    * for constructing fraud proofs of state transitions within these nested chains to be compatible with the tree structure while reframing the state transitions to be highly scalable
* a consensus mechanism which is dependent upon the root blockchain
    * to replicate the results of Nakamoto consensus incentives
* a bitmap UTXO commitment structure
    * for ensuring accurate state transitions off the root blockchain while minimizing mass-exit costs.


**Allowing for exits in data unavailability or other Byzantine behavior is one of the key design points in Plasma's operation**


### Design Stack and Smart Contracts

Plasma is not designed to reach assured finality raplidly, even though transactions are confirmed in the child chains raplidly, it requires it to be finalized on the underlying root blockchain.


#### Free Option Problem

In smart contracts, there is an issue of 'free option problem' whereby the receiver of a smart contract offer is needed to sign and broadcast the contract in order to enforce it - during that time the receiver of the contract may treat it as a free option and refuse to sign the contract if the activity does not interest them.

There is no guarantees of atomicity with the first and second signatures step for interactive protocol in blockchains.

A payment can split into many small payments.  This minimize the free option to the amount per split transaction.



### Multiparty Off-Chain State

There are two common issues in efforts to establish off-blockchain multiparty channels:

* the need to do synchronized state update amongest all participants when there need to be an update on the system, and therefore must be online.
* adding and removing participants in the channel require a large on-blockchain update, enumerating all participants which are added and removed.





