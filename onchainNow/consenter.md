
每个consenter需要实现

* type Consenter struct{}
	* 一个共识算法实现一个Consenter，
	* 其中的HandleChain接口拥有完成共识算法的实例化
* type Chain struct{}
* type ConsenterSupport struct{}


其中Consenter接口如下：

* HandleChain(support ConsenterSupport, metadata *cb.Metadata) (Chain, error)

	从metadata中取出所需要的信息，完成chain的实例化



