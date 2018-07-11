## RPC Message

```
type RPC struct {
	pb.RPC
	from		peer.ID
}

type pb.RPC struct {
	Subscriptions			[]*RPC_SubOpts
	Publish					[]*Message
	XXX_
}

message SubOpts {
	optional bool subscribe = 1
	optional string topicid = 2
}
```