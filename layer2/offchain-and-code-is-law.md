

## 链外扩容与Code is Law的思考

“链外扩容的思路，把比较复杂的计算任务放到链外去做，然后把数据的验证和真正需要全球共识的操作放到链上，剩下计算的部分全部放到链外。”

在一个区块链技术文章中看到上述的链外扩容思路，个人不敢苟同。也讨论一下个人对于区块链链外扩容的一些思考。

首先看上面的设计思路，其中依赖的核心点是，计算和验证的不对称性。如何实现计算和验证的不对称性，在理论上，是一个世界难题。如果理论上取得突破，可以实现通用业务逻辑的计算和验证的分离，那么带来的不仅仅是区块链技术的突破，更是安全计算的重大突破，在MPC和个人信息安全方面都有无限的应用场景。而在目前，只有密码学质数分解领域有现成的成果，比如ZKP，VDF等。在工程上，这个更加困难，远远超出了一般软件开发人员的能力范围。

第二点，在区块链上，code is law，这个code只限于是是区块链平台上运行的code，主要包括区块链上智能合约，虚拟机和账本状态管理相关的代码。按照上述设计，计算和验证分离的设计中，验证逻辑是位于区块链上，计算逻辑位于区块链下。因此对应的设计中，validation code is law, but production code is not law。因此，只要对其设计思路进一步思考，就不难发现，实际的设计要求是，production code和validation code在满足不对称复杂度的同时，实现逻辑的等价。从software engineering的角度看，这是一个充满风险的工程方向。

> TL;DR 做为一个工程师，提出设计同时必须首先考虑如何实现这个设计。
>
> 首先，我们看validation code。要完成一段代码必须要做的两个事情，一个是定义代码的功能，二是按照功能进行编码实现。对于validation code，由于validation code is law，其编码实现可以采用标准智能合约的开发方式完成。但是定义代码功能方面，必须对业务场景和业务规则作全面的约束和定义，如果我们将整个业务场景做为一个状态机，validation code是对整个状态机的完整定义。这对于规则简单场景，比如转账类，比较简单。但对于实际的业务场景，只能通过强加约束方式简化模型。但是，强加约束所需要付出的代价，就必须production code来付出了。
>
> 然后，我们看production code。production code的功能定义是很直接的，1. production code必须遵从the law defined by validation code，2. production code需要业务所需要的所有逻辑。production code不需要使用智能合约实现，其编码实现应该是更简单和方便的。但是，通常是事与愿违的，只要稍微了解software engineering的人，都知道软件的可靠性和逻辑复杂度是反比关系，软件的可测试性和逻辑的复杂度也是反比关系。由于validation code defined law，production code将必须100% compatible with its law。否则区块链平台将失去其根基。

第三点，依然是code is law，但是区块链没有律师，只有代码。而在实际的设计中，production code 不等于 validation code，违背了trustless platform的设计原则。production code处理用户交易，运行环境为服务提供方，validation code is law，运行于区块链上。按照信息安全的木桶原理，整个系统的安全是由服务提供方决定的，而不再trustless。

那么，是不是链外扩容的思路就是错误的呢？不是。实际上，在链外扩容方面早就有一个很好的设计方法，就是CounterFactual Instantiation。CounterFactual和上述思路有什么不同呢？用简单的一句话就是，在CounterFactual的设计中，validation code == production code。

关于CounterFactual的设计方法，且听下回分解。




