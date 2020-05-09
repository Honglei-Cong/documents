
# Stable and Probabilistic Leaderless BFT Consensus through Metastability

* traditional consensus algorithm
  * strong finality
  * BFT-style
  * focusing on log-ordering or state-machine-replication
* blockchain consensus algorith
  * probabilistic finality
  * Nakamoto / Avalanche
  * focusing on fork choosing


### Key Guarantees

* Safety : e-safety guarantee that is probabilistic
* Liveness
  * non-zero probability guarantee of ternimation within a bounded amount of time
  * time required for finality approaches `INF` as f approaches n/2

### Slush

* almost memoryless: retains no state between rounds other than its current color, maintains no history of interactions with other peers.
* every round involves sampling a small, constant-sized slice of the network at random
* make progress under any network configuration

```
procedure OnQuery(v, col')
    if col = nil then col = col'
    Respond(v, col)

procedure SlushLoop(u, col_0 in {R, B, nil})
    col := col_0
    for r in {1 ... m} do
        if col = nil the continue
        K = sample(N\u, k)
        P = [Query(v, col) for v in K]
        for col' in {R, B} do
            if P.cound(col') >= alpha then
                col = col'
    accept(col)
```

### Snowflake

Snowflake arguments Slush with a single counter that captures the strength of a node's conviction in its current color.  This per-node counter stores how many consecutive samples of the network by the node have all yielded the same color.  A node accpets the current color when its counter exceeds `beta`.

### Snowball

Snowball arguments Snowflake with confidence counters that capture the number of queries that have yielded a threshold result for their corresponding color.

### Multi-Value Consensus

Every node operates in three phases:
1. it gossips and collects proposals for `O(log n)` rounds, where each round lasts for the maximum message delay,
2. node stops collecting proposals, and instead gossips all new values for an additional `O(log n)` rounds,
3. each nodes samples the proposals it knowns of locally, checking for values that have an `alpha` majority, ordered deterministically, such as by hash values.
4. a node selects the first value by the order as its initial state when it starts the subsequent consensus protocol.


### AValanche

Avalanche consists of multiple single-decree Snowball instances instantiated as a multi-decree protocol that maintains a dynamic, append-only DAG of all known transactions.

**Transaction**
* names one or more parents, and form the edges of the DAG

```
procedure INIT
  {T} := {}   // the set of known transactions
  {Q} := {}   // the set of queried transactions

procedure OnGenerateTx(data)
  edges = {T' <- T : T' in ParentSelection(T)}
  T := Tx(data, edges)
  OnReceiveTx(T)

procedure OnReceiveTx (T)
  if T not in {T} then
    if P_T is {}, then
      P_T = {T}
      P_T.pref = T
      P_T.last = T
      P_T.cnt = 0
    else
      P_T = union(P_T, {T})
    {T} = {T} + T
    c_T = 0

confidence value: d_u(T) = sum of c_uT' if T' in {T_u} and T *<- T'

procedure AvalancheLoop
  while true do
    find T that satisfies T in {T} && T not in {Q}
    {K} = sample({N}\u, k)
    {P} = {every v in {K}, Query(v, T)}
    if |{P}| >= alpha then
      c_T = 1
      for T' in {T} : T' *<- T do     // T' is ancients of T
        if d(T') > d(P_T'.pref) then
          P_T'.pref = T'
        if T' != P_T'.last then
          P_T'.last = T'
          P_T'.cnt = 0
        else
          ++P_T'.cnt

    {Q} = {Q} + T

function IsPreferred(T)
  return T = P_T.pref

function IsStronglyPreferred(T)
  return all T' in {T}, T' *<- T: IsPreferred(T')

function IsAccepted(T)
  return P_T.cnt > beta_2
         or
         all T' in {T}, T' *<- T: IsAccepted(T'), and |P_T| = 1, and d(T) > beta_1

procedure OnQuery(j, T)
  OnReceiveTx(T)
  Respond(j, IsStronglyPreferred(T))

```


