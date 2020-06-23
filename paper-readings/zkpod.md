
zkPoD: practical decentralized system for data exchange


PoD : proof-of-delivery


In the past: placing their trust on fewer entities may suggest less risk: a smaller trusted base implies a more secure system.

Blockchain-based third party (BTP) has some significant advantages over its centralized counterpart.


### PoD Protocols

Three Parties: 
* Alice: data sender
* Bob: data receiver
* Julia: judge (trustless third party with visible internal states and predictable behaviors)

Three phases of PoD:
1. Init-phase 
    * All three parties setup with system parameters.
    * Alices publishes the authenticator of data (`\sigma`) by sending it to Julia such that everyone including outsiders can see it.
2. Deliver-phase
    * Alice encrypts the data `m` to get `m'` with a random one-time key `k` before she sends `m'` and the commitment of `k` to Bob.
    * Bob verifies:
        * encrypted data is consistent with the authenticator (`\sigma`)
        * data is indeed encrypted by `k`
    * Bob submits a `delivery receipt` containing the information of `k` if Bob accepts `m'`
3. Reveal-phase
    * After confirming the receipt, Alice reveals the key `k` to Julia to redeem the coins.
    * Julia verifies if the key matches the delivery receipt. (for trading, Julia may transfer coins, deposited by Bob, to Alice if the key is accepted)





