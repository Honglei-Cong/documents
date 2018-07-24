

## Multiplex Connection



Each peer has one ‘MConnection’ (multiplex connection) instance.

Each MConnection handles message transmission on multiple abastract communication channels.
Each channel has a globally unique byte id.

```
MConnection {
    cmn.BaseService

    conn                net.Conn
    bufConnReader       *bufio.Reader
    bufConnWriter       *bufio.Writer

    sendMonitor         *flow.Monitor
    recvMonitor         *flow.Monitor

    send                chan struct{}
    pong                chan struct{}

    channels            []*Channel
    ....
}
```

#### Channel

```
type Channel struct {
        conn            *MConnection
        desc            ChannelDescriptor
        sendQueue       chan []byte
        sendQueueSize   int32

        recving         []byte
        sending         []byte
}
```

#### BaseService接口

MConnection 实现了BaseService接口

* OnStart
* OnStop

#### MConnection Recv Routine

recvRoutine reads PacketMsgs and reconstructs the message using the channels 'recving' buffer.
After a whole message has been assembled, it is pushed to OnReceive()

Blocks depending on hwo the connection is throttled. Otherwise, it never blocks.

```
func (c *MConnection) recvRoutine() {
        for {
                c.recvMonitor.Limit(c._maxPacketMsgSize, atomic.LoadInt64(&c.config.RecvRate), true)

                var packet Packet
                _n, err := cdc.UnmarshalBinaryReader(c.bufConnReader, &packet, c._maxPacketMsgSize)
                c.recvMonitor.Update(_n)

                switch pkt := packet.(type) {
                case PacketPing:
                        c.pong <- struct{}{}            // notify to send pong msg

                case PacketPong:
                        c.pongTimeoutCh <- false        // notify pong received
                
                case PacketMsg:
                        channel, ok := c.channelsIdx[pkt.ChannelID]
                        if !ok || channel == nil {
                                // err handling
                        }
                        msgBytes, err := channel.recvPacketMsg(pkt)
                        if msgBytes != nil {
                                c.onReceive(pkt.ChannelID, msgBytes)    // reactor receives msg
                        }
                default:
                }
        }
}
```


#### MConnection Send Routine

sendRoutine polls for packets to send from channels.


```
func (c *MConnection) sendRoutine() {
        for {
        select {
        case <-c.flushTimer.Ch:
                c.flush()

        case <-c.chStatsTimer.Chan():
                for _, channel := range c.channels {
                        channel.updateStats()
                }

        case <-c.pingTimer.Chan():
                // send ping

        case timeout := <-c.pongTimeoutCh:
                c.stopPongTimer()


        case <-c.pong:
                cdc.MarshalBinaryWriter(c.bufConnWriter, PacketPong{})
                c.flush()

        case <-c.quit:
                break FOR_LOOP

        case <-c.send:
                c.sendSomePacketMsgs()      // send msgs
        }
        }
}
```


#### Send Message

queue a message to be sent to channel.

MConnection sendRoutine will iterate all channels, and send packages through the connection.


```
func (C *MConnection) Send(chID byte, msgBytes []byte) bool {
        channel, ok := C.channelsIdx[chID]
        channel.sendBytes(msgBytes)
}

func (ch *Channel) sendBytes(bytes []byte) bool {
        select {
        case ch.sendQueue <- bytes:
                atomic.AddInt32(&ch.sendQueueSize, 1)
                return true
        case <-time.After(defaultSendTimeout):
                return false
        }
}
```





