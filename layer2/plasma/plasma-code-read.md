
# Plasma Core

## Transaction

#### Fields 

* blknum1
* txindex1
* oindex1
* blknum2
* txindex2
* oindex2
* cur12
* newowner1
* amount1
* newowner2
* amount2
* sig1
* sig2

#### Property

* hash(self)
* merkel_hash(self)
* is_single_utxo(self)
* is_deposit_transaction(self)  : UTXO的输入为空
* sender1(self)
* sender2(self)
* encoded(self)
* sign1(self)
* sign2(self)


## Block

from rlp.Serializable

#### Fields

* TransactionSet
* Number
* Sig

#### Property

* hash (self)
* signer (self)
* merkle (self)
* root (self)
* is_deposit_block (self)
* encoded (self)
* sign (self, key)
* add_transaction (self, tx)

## Chain

#### Fields

* operator
* blocks = {}
* parent_queue = {}
* child_block_internval = 1000
* next_child_block
* next_deposit_block = 1

#### Funcs

* add_block(self, block)
* validate_transaction (self, tx, temp_spent={})
* get_block(self, blknum)
* get_transaction (self, utxo_id)
* mark_utxo_spent(self, utxo_id)
* \_apply\_transaction (self, tx)
* \_validate\_block (self, block)
* \_apply\_block (self, block)


# Plasma

## Root Chain

### Deployer

**Functions**

* deploy\_contract (self, contract_name, gas, args, consise)
* get\_contract\_at\_address (self, contract_name, address, consise)


### Root Chain Contract


**Events**

* Desposit (depsitor, depositBlock, token, amount)
* ExitStarted (exitor, utxoPos, token, amount)
* BlockSubmitted (root, timestmap)
* TokenAdded (token)

**Storage**

* const EXIT\_BOND = 1234567890
* const CHILD\_BLOCK\_INTERVAL = 1000
* operator
* currentChildBlock
* currentDepositBlock
* currentFeeExit
* mapping (uint256 => PlasmaBlock) plasmaBlocks
* mapping (uint256 => Exit) exits
* mapping (address => address) exitsQueues

```
struct Exit { owner, token, amount }
struct Plasmablock { root, timestamp }
```

**Functions**

* submitBlock (_root) 
    * operator将child chain block root提交到 root chain
    * 触发 BlockSubmitted (_root, block.timestamp)

* deposit ()
    * allow anyone to deposit funds into child chain
    * 构建新的deposit PlasmaBlock，添加到下一个depositBlock位置，质押 msg.value 到 child chain
    * 触发 Deposit (msg.sender, depositBlock, 0, msg.value)

* startDepositExit (depositPos, token, amount)
    * start an exit from a deposit
    * 找到质押的区块，即要退出的 utxo 的区块，msg.sender作为exitor，添加到 exit queue

* startFeeExit (token, amount)
    * allow the operator withdraw any allotted fees.
    * 以当前区块，添加到exit queue

* startExit (utxoPos, txBytes, proof, sigs)
    * start to exit a specific utxo
    * 找到utxo的区块，utxo的owner做完exitor，添加到exit queue

* challengeExit (cUtxoPos, eUtxoIndex, txBytes, proof, sigs, confirmationSig)
    * allows anyone to challenge an exiting transaction by submitting proof of double spend on the child chain.
    * cUtxoPos：做出challenge的utxo
    * eUtxoIndex：正在退出的utxo
    * 通过eUtxoIndex 检查exits[]中是否正在退出
        * 如果确认，将退出请求从exit queue中删除
    * msg.sender收到奖励

* getNextExit (token)
    * determine the next exit to be processed

* finalizeExits(token)
    * process any exits that have completed the challenge period

* getPlasmaBlock (blockNumber)
    * query the child chain

* getDepositBlock ()
    * determines the next deposit block number

* getExit (utxoPos)
    * returns information about an exit

* addExitToQueue (utxoPos, exitor, token, amount, createdAt)
    * add an exit to the exit queue
    * 按照utxo的创建时间为优先级，放到队列中
    * 触发 ExitStarted(msg.sender, utxoPos, token, amount)


## Child Chain


**Storage**

* operator
* root_chain
* chain = Chain(operator)
* current\_block


