

## Deposit in Plasma

To start real transactions on Plasma, tokens must be deposited firstly.

1. The tokens are sent into Plasma contract on the root blockchain.
   The tokens are recoverable within some set time period for a challenge/response.

2. The Plasma blockchain includes an incoming transaction proof.
   When this is included, the blockchain is committing to the fact that it will honor a withdrawal request.

3. Depositor signs a transaction on the child Plasma blockchain
   activating the transaction, which includes a commitment that they have seen the block with the chain's commitment in Phase2.


After this process, the chain has committed to the fact that they will handle these tokens and gave allocation so that withdrawls can be compactly proven.  With the third phase, the user is attesting to the fact that they can withdraw.


