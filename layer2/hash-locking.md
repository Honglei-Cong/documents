
## Hash-Locking

https://bitcointalk.org/index.php?topic=193281.msg2224949#msg2224949

#### Prepare

A picks a random number x, h = H(x)

#### Deploy Contract

A creates TX1: "Pay w BTC to <B's public key> if (x for H(x) known and signed by B) or (signed by A & B)"

A creates TX2: "Pay w BTC from TX1 to <A's public key>, locked 48 hours in the future, signed by A"

A sends TX2 to B

B signs TX2 and returns to A

#### Exchange

1) A submits TX1 to the network

B creates TX3: "Pay v alt-coins to <A-public-key> if (x for H(x) known and signed by A) or (signed by A & B)"

B creates TX4: "Pay v alt-coins from TX3 to <B's public key>, locked 24 hours in the future, signed by B"

B sends TX4 to A

A signs TX4 and sends back to B

2) B submits TX3 to the network

3) A spends TX3 giving x

4) B spends TX1 using x

## Explain

假设Alice想用自己的ADA换取Bob的BTC，具体应用流程如下：

1. Alice创建了一个随机密码s，并且算出该密码的哈希值h，即h = hash(s)。Alice将这个哈希值h发给Bob。

2. Alice和Bob共同通过智能合约将彼此的资产先后锁定（Alice先进行锁定，Bob再锁定），并且智能合约里实现了以下逻辑：
    * 条件1：如果任何人能在H小时内提供一个随机数值s'给智能合约，一旦合约验证了hash(s') == h（当s'等于原始密码s），那么Bob的BTC就自动转给Alice，否则超时后发还给Bob。
    * 条件2：如果任何人在2H小时内将原始密码s发给智能合约，则Alice的ADA将被自动转给Bob，否则转还给Alice。

上述的条件1是针对Alice来制定的，Alice为了拿到Bob的BTC，必然会在超时前将自己产生的随机密码s提供给智能合约，合约验证肯定会顺利通过，并且把Bob的资产转给Alice。

这一笔交易在成功的同时，Alice提供的原始密码s也被公开地广播并记录在了区块链上。此时，Bob可以拿着公开了的密码s，发给智能合约，依照条件2，便能获得Alice锁定在智能合约中的ADA，理论上Bob有1H到2H小时的充裕时间来完成操作（取决于Alice多快能完成她的操作）。

## Analysis

This is atomic (with timeout).  If the process is halted, it can be reversed no matter when it is stopped.

Before 1: Nothing public has been broadcast, so nothing happens

Between 1 & 2: A can use refund transaction after 48 hours to get his money back

Between 2 & 3: B can get refund after 24 hours.  A has 24 more hours to get his refund

After 3: Transaction is completed by 2

* A must spend his new coin within 24 hours or B can claim the refund and keep his coins
* B must spend his new coin within 48 hours or A can claim the refund and keep his coins

For safety, both should complete the process with lots of time until the deadlines.


