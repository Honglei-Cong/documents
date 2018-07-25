
## Deploy Guide

#### CA Setup

Start CA server

CA Server will generate its root Key.

CA Server supports hierarchical configuration.  
Intermediate CA can be setup to hidden Root CA from public.

#### Generate Certs from CA

Using fabric-ca-client to acquire cert from CA server.

```
fabric-ca-client enroll ...
```

#### Configurate Chain

All the following information should be configured in configtx.yaml

* how many participants
* participant ID
* Cert files
* participant security config
    * hash
    * crypto config
* participant nodes
    * host IP
    * host port
* orderer config
    * type
    * nodes config
    * block size


#### Creating Genesis Block

Using configtxgen to create genesis block and channel-create transaction.

Genesis block is for orderer to initialize chain.
Channel Creation Transaction is used for peer node to initialize channel.


#### Setup Orderer Network

Orderer network can be setup through docker compose.

All orderer nodes will boot up from genesis block, and initialize its block store.

Order will handle the configtx transaction, can setup new channel in orderer.
Of course, configtx transaction should be consensused on blockchain firstly.

#### Setup Peer Network

Peer network also can be setup through docker compose.

One of peers will initialize the channel-creation transaction, and channel will be created by orderers if the channel creation transaction is valid.

All other peers will get channel genesis block firstly through 'peer channel block get?'.  Then, peers join the channel with 'peer channel join' command.


#### Chaincode



