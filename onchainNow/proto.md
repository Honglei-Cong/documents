
plan to use protobuf for serialization.

### Consensus Message

```
// sent for every step taken in the ConsensusState, for every height/round/step transition
NewRoundStepMessage{
    Height        uint64
    Round        int
    Step            RoundStepType { NewHeight, NewRound, Propose, Prevote, PrevoteWait, Precommit, PrecommitWait, Commit }
    SecondsSinceStartTime       int
    LastCommitRound                int
}

// send when a block is committed (完成commit后的广播）
CommitStepMessage {
    Height                        uint64
    BlockPartsHeader     PartSetHeader {
                Total        int
                Hash        HexBytes
    }
    BlockParts                 BitArray
}

// sent when a new block is proposed （发出区块提案）
ProposalMessage {
    Proposal            *Proposal {
            Height                uint64
            Round                int
            Timestamp        Time
            BlockPartsHeader            PartSetHeader
            POLRound                        int
            POLBlockID                      BlockID                    (Proof-of-Lock round, votes from a previous round)
            Signature
    }
}

// sent when a previous proposal is re-proposed
ProposalPOLMessage {
    Height                             uint64
    ProposalPOLRound        int
    ProposalPOL                   BitArray
}

// sent when gossiping a piece of the proposed block
BlockPartMessage {
    Height                    uint64
    Round                    int
    Part                        *Part {
            Index                int
            Bytes               HexBytes
            Proof               SimpleProof
            hash                []byte
    }
}

// sent when voting for a proposal (or lack thereof)
VoteMessage {
    Vote            *types.Vote {
            ValidatorAddress        Address
            ValidatorIndex            int
            Height                         int64
            Round                         int
            Timestamp                 Time
            Type                            byte
            BlockID                       BlockID
            Signature
    }
}

// sent to indicate that a particular vote has been received
HasVoteMessage {
    Height                int64
    Round                int
    Type                   byte
    Index                  int
}

// sent to indicate that a given BlockID has seen +2/3 votes
VoteSetMaj23Message {
    Height            int64
    Round            int
    Type               byte
    Index              int
}

// sent to communicate the bit-array of votes seen for the blockID
VoteSetBitsMessage {
    Height            int64
    Round            int
    Type                bype
    BlockID           BlockID
    Votes               *BitArray
}

// sent to signal that a node is alive and waiting for transactions for a proposal
ProposalHeartbeatMessage {
    Heartbeat        *Heartbeat {
            ValidatorAddress        Address
            ValidatorIndex            int
            Height                        int64
            Round                        int
            Sequence                  int
            Signature
    }
}
```

### WAL Message

```
types.EventDataRoundState {
	Height 		int64
	Round 			int
	Step 			string
	RoundState 	interface{}
}

msgInfo {
	Msg			ConsensusMessage
	PeerID		p2p.ID
}

timeoutInfo {
	Duration	time.Duration
	Height		int64
	Round		int
	Step		cstypes.RoundStepType {NewHeight, NewRound, Propose, Prevote, PrevoteWait, Precommit, PrecommitWait, Commit}
}

EndHeightMessage {
	Height 	int64
}

```

### p2p.conn

```

PacketPing struct {}

PacketPong struct {}

PacketMsg struct {
	ChannelID		byte
	EOF				byte
	Bytes			[]byte
}

```

