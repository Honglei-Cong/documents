
## History of Casper Research

https://twitter.com/VitalikButerin/status/1029901077687066625

1. Today I am going to make a tweet storm explaining the history and state of Ethereum's Casper research, including the FFG vs CBC wars, the hybrid => full switch, the role of randomness, mechanism design issues, and more.

2. Ethereum proof of stake research began in Jan 2014 with Slasher (https://blog.ethereum.org/2014/01/15/slasher-a-punitive-proof-of-stake-algorithm/ …). Though the algorithm is highly suboptimal, it introduced some important ideas, most particularly the use of penalties to solve the nothing at stake problem (https://ethereum.stackexchange.com/questions/2402/what-exactly-is-the-nothing-at-stake-problem …).

3. That said, the penalties that I used were very small, only canceling out signing rewards. Vlad Zamfir joined in mid-2014, and he quickly moved toward requiring validators to put down *deposits*, much larger in size than rewards, that could be taken away for misbehavior.

4. Here's Vlad's retelling: [https://t.co/ODnHkHHeWS]

5. We spent much of late 2014 trying to deal with "long range attacks", where attackers withdraw their stake from deposits on the main chain, and use it to create an alternate "attack chain" with more signatures than the main chain, that they could fool clients into switching to.

6. If the attack chain diverges from the main chain at a fairly recent point in time, this is not a problem, because if validators sign two conflicting messages for the two conflicting chains this can be used as evidence to penalize them and take away their deposits.


7. But if the divergence happened long ago (hence, long range attack), attackers could withdraw their deposits, preventing penalties on either chain.


8. We eventually decided that long range attacks are unavoidable for pretty much the reasons PoW proponents say (eg. https://download.wpsoftware.net/bitcoin/pos.pdf ). However, we did not accept their conclusions.

9. We realized that we could deal with long range attacks by introducing an additional security assumption: that clients log on at least once every four months (and deposits take four months to withdraw), and clients simply refuse to revert further than that.

10. This was anathema to PoW proponents because it feels like a trust assumption: you need to get the blockchain from some trusted source when you sync for the first time.

11. But to us dirty subjectivists, it did not seem like a big deal; you need some trusted source to tell you what the consensus rules of the blockchain are in any case (and don't forget software updates), so the additional trust required by this PoS assumption is not large.

12. Here's Vlad's retelling: [https://t.co/0hZBSzHOOf]

13. Now that we settled on deposits and penalties, we had to decide what those deposits and penalties _are_. We knew that we wanted an "economic finality" property, where validators would sign on blocks in such a way that ...

14 ...once a block was "finalized", no _conflicting_ block could be finalized without a large portion of validators having to sign messages that conflict with their earlier messages in a way that the blockchain could detect, and hence penalize.

15. I went on a big long, and ultimately unproductive, tangent on a direction I called "consensus by bet": [https://t.co/dIdFX0HCpv]

16. Consensus by bet was an interesting construction where validators would make bets on which block would be finalized, and the bets themselves determined which chain the consensus would favor.

17. The theory was that PoW also has this property, as mining is a bet where if you bet on the right chain, you gain (reward - mining cost), and if you bet on the wrong chain, you lose the mining cost, except with PoS we could push the odds on the bets much higher.

18. The odds on validators' bets would start off low, but as validators saw each other getting more and more confident about a block, everyone's odds would rise exponentially, in parallel, until eventually they would bet their entire deposits on a block. This would be "finality".

19. In the meantime, Vlad started heavily researching mechanism design, particularly with an eye to making Casper more robust against oligopolies, and we also started looking at consensus algorithms inspired by traditional byzantine fault tolerance theory, such as Tendermint.

20. Vlad decided that traditional BFT was lame (he particularly disliked hard thresholds, like the 2/3 in PBFT and Tendermint), and he would try to effectively reinvent BFT theory from scratch, using an approach that he called "Correct by Construction" (CBC)

21. In Vlad's own words: [https://t.co/gVvK1S8mMa]

22. The correct-by-construction philosophy is very different from traditional BFT, in that "finality" is entirely subjective. In CBC philosophy, validators sign messages, and if they sign a message that conflicts with their earlier message...

23 ... they have to submit a "justification" proving that, in the relevant sense, the new thing they are voting for "has more support" than the old thing they were voting for, and so they have a right to switch to it.

24. To detect finality, clients look for patterns of messages that prove that the majority of validators is reliably voting for some block B in such a way that there is no way they can switch away from B without a large fraction of validators "illegally" switching their votes.

25. For example, if everyone votes for B, then everyone votes on a block that contains everyone's votes for B, that proves that they support B and are aware that everyone else supports B, and so they would have no legitimate cause for switching to something other than B.

26. I eventually gave up on consensus-by-bet because the approach seemed too fundamentally risky, and so I switched back to trying to understand how algorithms like PBFT work. It took a while, but after a few months I figured it out.

27. I managed to simplify PBFT (http://pmg.csail.mit.edu/papers/osdi99.pdf …) and translate it into the blockchain context, describing it as four "slashing conditions", rules that state what combinations of messages are self-contradictory and therefore illegal: [https://t.co/ch1rcaPLLl]

28. I defined a rule for determining when a block is finalized, and proved the key "safety" and "plausible liveness" properties: (i) if a block is finalized, then there is no way for a conflicting block to get finalized without >= 1/3 violating a slashing condition ...

29. ... (ii) if a block is finalized, 2/3 honest validators can always cooperate to finalize a new block. So the algorithm can neither "go back on its word" nor "get stuck" as long as > 2/3 are honest.

30. I eventually simplified the minimal slashing conditions down from four to two, and from there came Casper the Friendly Finality Gadget (FFG), which is designed to be usable as an overlay on top of any PoW or PoS or other blockchain to add finality guarantees.

31. Finality is a very significant advancement: once a block is finalized, it is secure regardless of network latency (unlike confirmations in PoW), and reverting the block requires >= 1/3 of validators to cheat in a way that's detectable and can be used to destroy their deposits

32. Hence, the cost of reverting finality can run into the billions of dollars. The Casper CBC and Casper FFG approaches both achieve this, though in technically different ways.

33. Note that Casper CBC and Casper FFG are *both* "overlays" that need to be applied on top of some existing fork choice rule, though the abstractions work in different ways.

34. In simplest terms, in Casper CBC the finality overlay adapts to the fork choice rule, whereas in Casper FFG the fork choice rule adapts to the finality overlay.

35. Vlad's initial preference for the fork choice rule was "latest message-driven GHOST", an adaptation of GHOST (https://eprint.iacr.org/2013/881.pdf ) to proof of stake, and my initial preference was to start off with hybrid PoS, using proof of work as the base fork choice rule.

36. In the initial version of Casper FFG, proof of work would "run" the chain block-by-block, and the proof of stake would follow close behind to finalize blocks. Casper CBC was full proof of stake from the start.

37. At the same time, Vlad and I were both coming up with our own respective schools of thought on the theory of consensus *incentivization*.

38. Here, a very important distinction is between *uniquely attributable faults*, where you can tell who was responsible and so can penalize them, and *non-uniquely attributable faults*, where one of multiple parties could have caused the fault.

39. The classic case of a non-uniquely-attributable fault is going offline vs censorship, also called "speaker-listener fault equivalence".

40. Penalizing uniquely attributable faults (eg. Casper FFG slashing conditions) is easy. Penalizing non-unquely-attributable faults is hard.

41. What if you can't tell if blocks stopped finalizing because a minority went offline or because a majority is censoring the minority?

42. There are basically 3 schools of thought on this issue:

* Penalize both sides a little
* Penalize both sides hard (Vlad's preference)
* Split the chain into two, penalize one side on each chain, and let the market decide which chain is more valuable (my preference).

43. Or, in my words: [https://t.co/tUSwmLT882]

44. In November 2017, the Casper FFG slashing conditions, plus my ideas for solving "the 1/3 go offline" problem through a "quadratic leak" mechanism, became a paper: https://arxiv.org/abs/1710.09437 

45. Of course, I was well aware that appealing to the social layer to solve 51% attacks was not a very nice thing to do, so I started looking for ways to at least allow online clients to *automatically* detect which chain is "legitimate" and which is the "attack" in real time.

46. Here is one of my earlier ideas: https://ethresear.ch/t/censorship-rejection-through-suspicion-scores/305 … . It was something, but still suboptimal; unless network latency was exactly zero, there was only a guarantee that clients' suspicion scores would differ by at most delta, not that clients would fully agree. [https://t.co/Kem7zph3ZL]

47. In the meantime, my main criticism of Vlad's model had to do with "discouragement attacks", where attackers could credibly threaten to make a 51% attack that causes everyone to lose money, thereby driving everyone else to drop out, thus dominating the chain at near-zero cost

48. Vlad (along with Georgios Piliouras) started doing economic modeling to estimate the actual cost of such an attack under his model.

49. It's worth noting here that all of these issues are not unique to proof of stake. In fact, in proof of work, people tend to simply give up and assume preventing 51% attacks is outright impossible, and a 51% attack is a doomsday that must be prevented at all costs.

50. But, as is the Ethereum tradition, Vlad and I were both unaware that the word "ambitious" can be anything but a compliment, and kept on working on our separate approaches to disincentivizing, mitigating and recovering from 51% attacks.

51. In early 2018, Vlad's work on CBC started to move forward quickly, with great progess on safety proofs. For the state of progress in March 2018, see this epic two-hour presentation: [https://t.co/J3bBJVHm0T]

52. In the meantime, Casper FFG was making huge progress. A decision to implement it as a contract that would be published to the Ethereum blockchain made development easy. On Dec 31, 2017 at 23:40, we released a testnet written in python: [https://t.co/QbJT89Au7N]

53. Unfortunately, development of FFG then slowed down. The decision to implement FFG as a contract made some things easier, but it made other things harder, and it also meant that the eventual switch from EVM to EWASM, and single-chain Casper to sharded Casper, would be harder.

54. In addition, the team's work was being split between "main chain Casper" and "shard chain Casper" and it was clear there was enormous unneeded duplication of effort going on between the Casper and sharding teams.

55. In June 2018, we made the fateful decision to scrap "hybrid Casper FFG as a contract", and instead pursue full Casper as an independent chain, designed in such a way that integrating sharding would be much easier.

56. The switch to full proof of stake led me to start thinking much harder about proof of stake fork choice rules.

57. Casper FFG (and CBC) both require the *entire* validator set to vote in every "epoch" to finalize blocks, meaning there would be tens to hundreds of signatures coming in every second. BLS signature aggregation makes this practical in terms of computational overhead...

58. ...but I wanted to try to take advantage of all of these extra signatures to make the chain much more "stable", getting "100 confirmations" worth of security within a few seconds.

59. Here were my initial attempts: [https://t.co/PHag5oNMxw]

60. However, all of these approaches to the fork choice rule had a weakness: they split up validators into "attesters" and "proposers", and the proposers, being the key drivers of block production, had outsized power.

61. This was undesirable, primarily because it required us to have a strong source of on-chain random number generation to fairly pick the proposers. And on-chain randomness is *hard*, with simple approaches like RANDAO looking more and more problematic (https://ethresear.ch/t/randao-beacon-exploitability-analysis-round-2/1980 …). [https://t.co/ZaW2itzq2o]

62. Justin Drake and I went off to solve this problem in two ways, Justin by using verifiable delay functions which have a deterministic and verifiable output, but take a large amount of unparallelizable sequential time to compute, making manipulation ahead of time impossible...

63. ... and myself by making a major concession to the Cult of Vlad™, using GHOST-based fork choice rules to greatly reduce the dependence on proposers, allowing the chain to grow uninterrupted even if >90% of proposers are malicious, as long as >50% of attesters are friendly.

64. Vlad was very happy, though not fully: he preferred a version of GHOST based on validators' *latest messages*, whereas I preferred a version based on *immediate* messages: [https://t.co/Q8rBq6bcoC]

65. Around this time I also managed to come up with a way to "pipeline" Casper FFG, reducing time-to-finality from 2.5 epochs to the theoretically optimal 2 epochs: [https://t.co/0eutcy7Rhp]

66. I was very happy that the RPJ fork choice rule (which I have since renamed "immediate message-driven GHOST") is nicely compatible with Casper FFG in a way that most others are not...

67. ...and that it has a very important "stability" property: that the fork choice is a good prediction of the future fork choice. This seems obvious, but is very easy to accidentally make fork choice rules that do *not* have this property.

68. The most recent development of all is a result that latest message driven GHOST may, due to a technicality, only give 25% fault tolerance within two rounds, but immediate driven message GHOST (with FFG or CBC) still gives the full 33% (no writeup yet)

69. The main tradeoff between FFG and CBC is that CBC seems to have nicer theoretical properties, but FFG seems to be easier to implement.

70. In the meantime, a *lot* of progress on verifiable delay functions has been made: [https://t.co/XsRUmZH2ap]

71. Also, I recently decided to look into Leslie Lamport's old 1982 paper, where he had a consensus algorithm that has 99% fault tolerance if you add the assumption that all nodes, including observers, are online with low network latency: [https://t.co/PnDBzRRT9w]

72. The network latency assumptions arguably make this unsuitable as a primary consensus algorithm. However, there is one use case where it works *really* well: as a substitute for suspicion scores for 51% censorship detection.

73. Basically, if a 51% coalition starts censoring blocks, other validators and clients can detect that this is happening, and use the 99% fault tolerant consensus to agree that this is happening, and coordinate a minority fork.

74. The long-run goal of this research is to reduce reliance on the social layer as much as possible, and maximizing the cost of destabilizing the chain enough so that reverting to the social layer is necessary.

75. What's left now? On the FFG side, formal proofs, refinements to the specification, and ongoing progress on implementation (already started by >=3 teams!), with an eye to safe and speedy deployment. On the CBC side, much of the same. Onward and upward!


