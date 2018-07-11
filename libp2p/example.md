## Examples


### example1

```
sk, pk = crypto.GenerateEd25519Key(rand.Reader)
pid = peer.IDFromPublicKey(pk)
pstore = pstore.NewPeerStore()
pstore.AddPrivKey(pid, sk)
pstore.AddPubKey(pid, pk)
network := swarm.NewNetwork ({maddr}, pid, ps, nil)
host := bhost.New(network)

```

一个host如果有多个身份可以连接到多个network，一个network表示一个子网？

### exmaple2

* create basic host

```
transports := msmux.NewBlankTransport()
transports.AddTransport(“/yamux/1.0.0”, yamux.DefaultTransport)

swarm = swarm.NewSwarmWithProtector(pid, peerstore, transports)
basichost = bhost.New(swarm)
```

* create stream

```
host.SetStreamHandler(“/echo/1.0.0”, func (s net.Stream) {
	buf := bufio.NewReader(s)
	str = buf.ReadString(“\n”)
	s.Write([]byte(str)
	s.close()
})
```

* send data

```
stream = host.NewStream(peerid, “/echo/1.0.0”)
stream.Write([]byte(“hello world”)
```
