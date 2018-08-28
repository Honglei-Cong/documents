## Fabric Configuration

由几个大模块组成

* organizations
* capabilities
* Application
* Orderer
* Channel
* Profiles

### Organizations

定义各个Org，比如OrdererOrg, Org1, Org2等

每个Org的定义:

* Name
* ID
* MSPDir
* Policies : (readers, writers, admins)


### Capablities

定义各个模块的版本能力。

* Global
* Orderer
* Application


### Application

defines the values to encode into a config transaciton or genesis block for application related parameters.

* Organizations
    * list of orgs which are defined as participants on the application side of the network
* Policies
    * defines the set of policies at this level of the config tree
    * Readers, Writers, Admins
* Capabilities

### Orderer

定义orderer相关的参数

* OrdererType
* Addresses
* BatchTimeout
* BatchSize
    * MaxMessageCount
    * AbsoluteMaxBytes
    * PreferredMaxBytes
* Kafka
    * Brokers
* Organizations
* Policies
* Capabilities

### Channel

* Policies
    * Readers: who may invoke the 'deliver' api
    * Writers: who may invoke the 'broadcast' api
    * Admins: by default, who may modify elements at this config level
* Capabilities





