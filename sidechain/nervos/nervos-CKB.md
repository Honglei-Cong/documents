
## Nervos CKB

from Nervos whitepaper

### Current limitations
* scalability problems
    * Turing-complete VM, difficult for full nodes to determine dependencies between transactions.  This makes it difficult for nodes to process transactions in parallel.
* indeterministic state transition
    * contract state is updated by contract code, depends on the execution context.  User can not determine the exact execution result when they send the transaction.
* mono-contract
    * tightly coupled computation and storage.

###Components
* Computation
    * Generator (state generation)
    * Validator (state validation)
* Storage
    * Cell (storage)
    * Type (date type)
* Identity


### Generators
* implemented application logic
* runs on the client side to generate new states

### Validator
* state validation logic
* consensus node first authenticate the submitter of the transaction, 
* then validate new states in new transaction with Validators
* once a new block is generated and received, new states in block are committed into the CKB.

### Cell Data Source
* independently verifiable (trust built-in)
* data endorsed by identities in the system

#### storage unit in CKB
* data to be stored
* validator logic of the data
* data owner

### Cells

* type
* capacity : byte limit of data that can be stored in the cell
* data: the actual binary data store in the cell
* owner_lock: lock script to represent the ownership of the cell.
* data_lock: lock script to represent the user with right to write the cell.

Cell is an immutable data unit, and it can not be modified after creation.  Cell ‘updates’ are essentially creating new cells with the same ownership.  User create cells with new data through transactions, and invalidate old cells at the same time.

### Type

To define a new cell type, we must specify two essential components:

* Data Schema : the data structure of cells
* Validator : defines the validating rules of cells

### Transaction

* deps: dependent cell set, provides read-only data that validator needs.
* inputs: input cell set, includes cells to be transferred and/or updated.
* outputs: includes all newly created P1 cells.


### Nervos Nodes
* Archive Nodes
    * full nodes in the CKB network
    * validate new blocks and transactions, relay blocks and transactions
    * keep all historical transactions on disk
* Consensus Nodes
    * listen to new transactions, package them to blocks, and achieve consensus on new blocks
* Light Clients
    * store only very limited data, and can run on desktop computers or mobile devices


