
# "Some Simple Economics of the Blockchain"

For an exchange to be executed, key attributes of a transaction need to be verified by parties involved.
The need for intermediation increases as markets scale in size and reach both geographically and in terms of the number of participants involved.
When the cost of successfully verifying the relative attributes of a transaction outweights the benefits from the exchange, the market falls apart.

Common solutions to these problems involve relying on an intermediary for third-party verification, to maintain a reputation system, to force additional disclosures on the seller, to enforce contract clauses designed to generate a separating equilibrium (e.g. warranties), and to perform monitoring.

(随着交易scale的增大，不得不为exchange增加更多中间商，进行交易信息的验证。为保证中间商的信誉问题，有不得不增加监管部门。)
(另外，众多第三方的引入，同时引入了交易双方的信息泄漏，中间商将能够利用历史信息获利。)

Digitization, in particular, has pushed verification costs for many types of transactions close to zero. Blockchain technology has the potential to complete this process by allowing for the first time for distributed, costless verification.
(数字化极大降低了交易手续费，降低的这些手续费大大地扩大了中间商市场。区块链+数字化，可以进一步真正降低端到端的交易手续费。)

Blockchain:
* ability to perform an audit at zero cost, enables distributed, costless verification
* rules of the audit can be decided ex-ante, reducing the risk of a conflict of interest arising ex-post between the entity in charge of the audit and either side of the market
* privacy enhancing

(随着verification cost的降低，交易粒度将更小，之前微小价值的事物都将获得全新的market)

### Market Design and the Blockchain

**PoW的Market design**
In proof-of-work systems, “mining” does not serve the purpose of verifying transactions (this activity is fairly light computationally), but of building a credible commitment against an attack.
As a result, in proof-of-work systems, a blockchain is only as secure as the amount of computing power dedicated to mining it.
This generates economies of scale and a positive feedback loop between network effects and security: 
as more participants use a cryptocurrency, the value of the underlying token increases (because the currency becomes more useful), which in turn attracts more miners (due to higher rewards), ultimately increasing the security of the ledger.

(同时，PoW共识的效率也隐含了自身的Market Design。BTC采用10分钟出块，LTC采用2.5分钟，因此LTC更适合于相对小额的交易。每个区块链不同的market design，在某种意义下也意味着这个区块链适合哪种类型的交易。)

交易的属性：
* size of a transaciton
* its attributes
* its functionality
* related degree of security and privacy needed

When combined with privacy-enhancing measures, this can solve the trade-off between users’ desire for customized product experiences (e.g. when using a virtual assistant like Siri), and the need to protect their private information (e.g. the queries sent to the service).
If the sensitive data is stored on a blockchain, users can retain control of their data and license it out as needed over time.


Marketing:
* exchange
* address market failures
* monitor market participants as a substantially lower cost


### Applications

* linking a hardware device (e.g. an Internet of Things device, a solar panel) to a cryptocurrency. If the hardware device is secure and cannot be tampered with, then the information it collects can act as the trusted oracle in a digital transaction.

* a car that can read information from a blockchain and use public-key cryptography to authenticate its user

* Internet of Things (IoT) devices and robots, when combined with a cryptocurrency, can seam- lessly earn, barter or exchange resources with other devices on the same network



# "Market Design for a Blockchain-Based Financial System"

Under proof of stake, user's utility in equilibrium is generally above the fork threat level because custodians can use relational contracts to incentive a higher quality of service.
Relational contracts under proof of stake rely only on local institutions —- but combining them with cryptography can create a platform for formal global contracts.


### Tranditional Financial Systems

Traditional financial systems maintain trust through formal and relational contracts supported by enforcement mechanisms such as courts.

* increase efficiency
* higher transaction fee
* barriers to entry
* potentially suboptimal level of innovation
* reduce the resillience of a financial system to failure
* interference by third parties


### Definitions

nodes (process and verify transactions) == custodians (which hold users' asset)

relational contracts (between nodes and custodians) can secure the system and hence reduce the risk of catastrophic failure.

#### Proof of Work

In PoW, influence of a node over transactions is a function of that node's relative computational power.

Relational contracts between nodes and custodians are not enforceable because it is not possible to selectively exclude misbehaving nodes from the system.

In long-run equilibrium, the threat of a fork ensures a baseline level of utility for the users.

#### Proof of Stake

Under PoS, the influence of a node over transactions is a function of the relative stake that the node controls either directly or through delegation.

The use of relational contracts may also enable faster coordination around governance decisions.
At the same time, however, the use of relational contracts under proof of stake could make the system more vulnerable to interference by third parties.

Nodes’ identities and off-network reputations can play a role in sustaining equilibrium under proof of stake.

The presence of relational contracts imparts custodians with a more prominent role in the ecosystem. If trust in custodians becomes an important issue, then the market for custody might become concentrated; this concentration might be undesirable per se, and it might lead to another vector for third party interference.
(PoS依赖于custodian的可信性，从而导致custody market的进一步集中化。)

#### Conclusion

With weak relational contracts or substantial concerns about outside interference, proof of work designs may be preferable. 
In contrast, when some regions have local institutions that are reliable enough to make delegation feasible, proof of stake designs can lead to efficiency gains and improvements in governance.

If successful, proof of stake designs would then use cryptography to leverage these local institutions to create a platform for formal global contracts—providing spillover benefits to regions with weaker local institutions.
