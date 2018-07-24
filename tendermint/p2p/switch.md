
## Switch

Switch handles peer connections and exposes an API to receive incoming messages on 'Reactors'.

Each 'Reactor' is responsible for handling incoming messages of one or more 'Channels'.  
So while sending outgoing messages is typically performed on the peer, incoming messages are received
on the reactor.

```
type Switch struct {
        BaseService

        config          *config.P2PConfig
        listeners       []Listener
        reactors        map[string]Reactor

        chDescs         []*conn.ChannelDescriptor
        reactorsByCh    map[byte]Reactor

        peers           *PeerSet
        dialing         *CMap
        reconnecting    *CMap

        nodeInfo        NodeInfo
        nodeKey         *NodeKey
        addrBook        AddrBook

        mConfig         conn.MConnConfig
}
```


#### Listen Routines

start listeners

```
func (sw *Switch) listenrRoutine(l Listener) {
        for {
                inConn, ok := <-l.Connections()

                // ignore connection if we already have enough
                if maxPeers <= sw.peers.Size() {
                        continue
                }

                // new inbound connections
                sw.addInboundPeerWithConfig(inConn, sw.Config)
        }
}
```

On new peer connected, addPeer() performs the Tendermint P2P handshake with a peer that 
already has a SecretConnection.  If all goes well, it starts the peer and adds it to the switch.


#### Broadcast

Broadcast runs a go routine for each attempted send, which will block trying to send for 
defaultSendTimeoutSeconds.  Returns a channel which receives sucess values for each attempted send.
Channel will be closed once msg bytes are sent to all peers (or time out).


```
func (sw *Switch) Broadcast(chID byte, msgBytes []byte) chan bool {
        successChan := make(chan bool, len(sw.peers.List()))

        var wg sync.WaitGroup
        for _, peer := range sw.peers.List() {
                wg.Add(1)
                go func(peer Peer) {
                        defer wg.Done()
                        success := peer.Send(chID, msgBytes)
                        successChan <- success
                }
        }
        go func() {
                wg.Wait()
                close(successChan)
        }

        return successChan
}
```





