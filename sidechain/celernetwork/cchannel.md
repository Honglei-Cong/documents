## cChannel

To support arbitrary conditional dependency between on-chain verifiable states.

**Design Goals**

* achieve fast, flexible and trust-free off-chain interactions
* in most cases, off-chain state transitions will stay off-chain until final resolution
* design data structure and object interaction logic that works for different blockchains

**State**
state of a channel

**State-Proof**
sp = {delta_s, seq, merkle_root, sigs}

* accumulative state updates up until now
* the sequence number for the state proof
* root of the merkle tree of all pending condition groups
* signatures from all parties on this state proof
Condition: cond = {timeout, \*IsFinalized(args), *QueryResult(args)}
* timeout after which the condition expires
* bool function IsFinalized, to check whether the condition has been resolved and settled before the condition timeout
* QueryResult, return arbitrary bytes as the condition’s resolving result.
Condition Group: cond_group = {A, ResolveGroup(cond_results)}
* A: a set of conditions contained
* ResolveGroup: takes the resolving results of all conditions as inputs and returns a state update delta_s.

### State Channel

* p : set of participants
* s0 : onchain based state for this channel
* sp : most updated known state proof for the channel
* s: the updated channel state after state proof sp is fully settled
* t : settlement timeout
* F : contains a set standard functions that should be implemented by every state channel.
    * ResolveStateProof (sp, cond_groups)
    * GetUpdateState (sp, s0)
    * UpdateState (s)
    * IntendSettle (new_sp)
    * ConfirmSettle (sp)
    * IsFinalized (args)
    * CloseStateChannel (s)

