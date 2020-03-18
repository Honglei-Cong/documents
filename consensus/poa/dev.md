
### start

1. NewPoAServer
2. initialize
3. OpenBlockStore
4. newBlockPool
5. LoadChainConfig
6. for all peers : peerPool.addPeer
7. start
   1. syncer run
   2. statemgr run
   3. msgsend loop
   4. timer loop
   5. action loop
   6. goroutine : processMsgEvent
8. stateMgr <- ConfigLoaded
   * peer heartbeat sync block height with other nodes
   * block height synced, change state to sync-ready
   * sync-ready timeout, start new round
9. start
   1.  start peer ticker
   2.  for all peers : run(peer pubkey)

### Peer routine

1. send heartbeat
2. wait min connections
3. routine
   1. msg = get msg from peer-msg-channel
   2. decode consensus msg
   3. verify msg
   4. server.onConsensusMsg(fromPeer, msg, msg_hash)
   5. switch msg type

### procedures

1. if no higher msg available
   1. if leader of some round
      1. if proposal-justfiy available, 
         1. if txpool not empty, make proposal
         2. else wait txpool
      2. else wait proposal-justify
   2. else
      1. wait new-view timeout
2. else
   1. for all forks of next height
      1. chose heaviest fork
      2. if prepare-QC availble: vote for commit
      3. else if proprep-QC availbe: vote for prepare
      4. else if proactor: vote for proprepare


### Vote Msg Processing

1. verify message 
   1. verify signatures in msg (in msg verification)
   2. verify rounds in vote is in a chain
   3. verify no new-view-justify from peer outside the chain
   4. verify no other commit-msg from peer outside the chain
   5. verify proposal msg can only be at tail
2. for each round in msg:
   1. process new-view-justify
      1. if commit-locked on other fork:
         1. **pending for new-view commit**
      2. remove all other forks
         1. if leader of new view
            1. if txpool empty: start txpool timer
            2. else: make proposal, add proposal to pending-response
      3. start new-view timeout
   2. process exec-merkle-msg
      1. update block tree with new commit msg
      2. get updated block-chain-range
      3. for round in range
         1. if post-sealed on fork
            1. set post-sealed, send block-postseal-action
         2. else if sealed on fork
            1. set sealed, send block-seal action
         3. else if **locked on another fork**
            1. ignore
         4. else if prepare-qc on fork
            1. set committed, locked on fork, send block-commit msg
         5. else if proactor of fork, not prepare for sibling
            1. set prepared, send block-prepare msg
         6. else if leader of fork
            1. send block-propose msg
   3. process commit
   4. process prepare
   5. process proposal
      1. find blocknode
      2. if locked on other forks, return
      3. if prepared, return
      4. if proprepared
         1. make prepare
      5. else if is proactor
         1. if prev-blocknode is proprepared
            1. make prepare
3. response
   1. build response msg
   2. broadcast

### ViewChange Msg Processing

1. verify message
2. add message to block tree
3. if QC reached
   1. update block tree



### timer event

timer for each round, update timer on receiving any message in the round

* New-View timeout


### action

* seal-block (proposal)
* postseal-block (proposal)
* make-proposal (height, view)

### Chain Config

chain config update after 1+ consensus epoch.

chain config:
* Version
* View
* N
* C
* BlockMsgDelay
* HashMsgDelay
* PeerHandshakeTimeout
* Peers
  * Index
  * ID (pubkey)
* PosTable
* NextBlockChangeView

### Block Participant Config

Block Participant is committee, updated per block.

BlockParticipantConfig
* BlockNum
* LastVRF
* ChainConfig
* Leader
* proactors


### TODO

1. check msg pool when new msg processed
2. wait txpool when making new proposal
3. 
