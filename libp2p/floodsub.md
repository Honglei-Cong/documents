## Floodsub

```
type PubSubRouter interface {		// FloodSubRouter实现此接口
	Protocols() 		[]protocol.ID
	Attach(*PubSub)
	AddPeer(peer.ID, protocol.ID)
	RemotePeer(peer.ID)
	HandleRPC(*RPC)
	Publish(peer.ID, *pb.Message)		 // 将消息publish出去
}

type PubSub struct {
	counter

	host			host.Host				// 
	rt				PubSubRouter			// 

	incoming		chan *RPC				// 其它节点的消息
	publish			chan *Message			// 发送给我们的消息
	addSub			chan *addSubReq		// control channel，添加sub
	getTopics		chan *topicReq			// 请求topic列表
	getPeers		chan *listPeerReq		// 请求peer列表
	cancelCh		chan *Subscription		// 请求取消订阅
	newPeers		chan inet.Stream		// notify通道，通知新的peer加入
	peerDead		chan peer.ID			// notify通道，通知peer dead

	myTopics		map[string]map[*Subscription]struct{}		// 我们订阅的topic
	topics			map[string]map[peer.ID]struct{}		// 哪个peer sub了哪个topic

	sendMsg		chan *sendReq			// 要发送的消息，发到这个channel
	addVal			chan *addValReq		// handle validator registration requests
	rmVal			chan *rmValReq
	topicVals		map[string]*topicVal	// topic validators
	validateThrottle	chan struct{}			// limit the number of active validate thrd

	peers			map[peer.ID]chan *RPC	// 给每个peer发消息的RPC 通道
	seenMessages	*timecache.TimeCache	// 记录已经收到的消息
}

func NewPubSub(ctx, host.Host, rt PubSubRouter, …Option) (*PubSub, error)
func (p *PubSub) processLoop(ctx context.Context)

```

### HandleNewStream

在每个新stream，启动routine，读取RPC消息，将其转发到incoming通道，processLoop中将监测incoming通道，调用handleIncomingRPC处理RPC消息。在handleIncomingRPC中，如果需要，将实现消息转发。

host的stream handler为HandleNewStream，负责一个新的stream的消息处理


### PushMsg

实现了Publish的接口

```
func (p *PubSub) pushMsg(vals []*topicVal, src peer.ID, msg *Message) {
	if len(vals) > 0 {
		// 需要validate，走验证流程
		select {
		case p.validateThrottle <- struct{}{}:			// 通过channel实现流控
			go func() {
				p.validate(vals, src, msg)			// 验证msg，结束后maybePublish
				<-p.validateThrottle
			}
		}
		return
	}
	p.maybePublishMessage(src, msg.Message)		// 直接发送
｝
```

PushMsg有两个入口：

1. publish通道的请求，
2. 每个incoming的RPC请求，如果RPC请求中有需要publish的消息
  * 目前只发现ipfs publish命令使用了publish


### maybePublishMessage

```
// pubsub的主要路由逻辑
func (p *PubSub) maybePublishMessage(from peer.ID, pmsg *pb.Message) {
	id := msgID(msg)
	if p.seenMessage(id) {
		return
	}
	p.markSeen(id)
	p.notifySubs(pmsg)
	p.rt.Publish(from, pmsg)
}
```





