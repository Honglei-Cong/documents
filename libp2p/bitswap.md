
## bitswap on p2p


```
func NewFromIpfsHost(host host.Host, r routing.ContentRouting) BitSwapNetwork {
	bitswapNetwork := impl {
		host: host,
		routing: r,
	}

	host.SetStreamHandler(ProtocolBitswap, bitswapNetwork.handleNewStream)
	...
	
	host.Network().Notify((*netNotifiee)(&bitswapNetwork))
	
	return &bitswapNetwork
}
```

host.SetStreamHandler 注册stream处理函数


```
func (bsnet *impl) handleNewStream(s inet.Stream) {
	defer s.Close()
	if bsnet.receiver == nil {
		s.Reset()
		return
	}
	
	reader := ggio.NewDelimitedReader(s, inet.MessageSizeMax)
	for {
		received, err = bsmsg.FromPBReader(reader)
		if err != nil {
			if err != io.EOF {
				s.Reset()
				go bsnet.receiver.ReceiveError(err)
			}
			return
		}
		
		p := s.Conn().RemotePeer()
		bsnet.receiver.ReceiveMessage(ctx, p, received)
	}
}
```

在receiver.ReceiveMessage中完成消息处理



