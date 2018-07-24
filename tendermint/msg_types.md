
## Tendermint Message Type

#### summary

* NewRoundStepMessage
* CommitStepMessage
* ProposalMessage
* ProposalPOLMessage
* BlockPartMessage
* VoteMessage
* HasVoteMessage
* VoteSetMaj23Message
* VoteSetBitsMessage
* ProposalHeartbeatMessage



#### NewRoundStepMessage

sent for every step taken in the ConsensusState, for every height/round/step transition

```
NewRoundStepMessage {
    Height                  uint64
    Round                   int
    Step                    RoundStepType { NewHeight, NewRound, Propose, Prevote, PrevoteWait, Precommit, PrecommitWait, Commit }
    SecondsSinceStartTime   int
    LastCommitRound         int
}
```

#### CommitStepMessage

send when a block is committed (完成commit后的广播）

```
CommitStepMessage {
    Height               uint64
    BlockPartsHeader     PartSetHeader {
                Total    int
                Hash     HexBytes
    }
    BlockParts           BitArray
}
```

#### ProposalMessage

sent when a new block is proposed （发出区块提案）

```
ProposalMessage {
    Proposal            *Proposal {
            Height            uint64
            Round             int
            Timestamp         Time
            BlockPartsHeader  PartSetHeader
            POLRound          int
            POLBlockID        BlockID      // (Proof-of-Lock round, votes from a previous round)
            Signature
    }
}
```

#### ProposalPOLMessage

sent when a previous proposal is re-proposed

```
ProposalPOLMessage {
    Height           uint64
    ProposalPOLRound int
    ProposalPOL      BitArray
}
```


#### BlockPartMessage

sent when gossiping a piece of the proposed block

```
BlockPartMessage {
    Height         uint64
    Round          int
    Part           *Part {
            Index  int
            Bytes  HexBytes
            Proof  SimpleProof
            hash   []byte
    }
}
```

#### VoteMessage

sent when voting for a proposal (or lack thereof)

```
VoteMessage {
    Vote    *types.Vote {
            ValidatorAddress  Address
            ValidatorIndex    int
            Height            int64
            Round             int
            Timestamp         Time
            Type              byte
            BlockID           BlockID
            Signature
    }
}
```

#### HasVoteMessage

sent to indicate that a particular vote has been received

```
HasVoteMessage {
    Height int64
    Round  int
    Type   byte
    Index  int
}
```

#### VoteSetMaj23Message

sent to indicate that a given BlockID has seen +2/3 votes

```
VoteSetMaj23Message {
    Height int64
    Round  int
    Type   byte
    Index  int
}
```

#### VoteSetBitsMessage

sent to communicate the bit-array of votes seen for the blockID

```
VoteSetBitsMessage {
    Height  int64
    Round   int
    Type    bype
    BlockID BlockID
    Votes   *BitArray
}
```

#### ProposalHeartbeatMessage

sent to signal that a node is alive and waiting for transactions for a proposal

```
ProposalHeartbeatMessage {
    Heartbeat        *Heartbeat {
            ValidatorAddress Address
            ValidatorIndex   int
            Height           int64
            Round            int
            Sequence         int
            Signature
    }
}
```





