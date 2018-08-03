## State Channel Applications

https://medium.com/statechannels/state-channel-applications-1f170e7d542e

### Current Dapp Design Pattern

* A public API - public or external functions on their contract
* Authorization logic - usually by checking msg.sender
* Resolution logic - to decide how funds are distributed


### Design a State Channel Application

* Not signing Eth transactions anymore, you're signing generic state objects
* Authorization logic requires each user of the state channel to sign an object for each update and verify these signatures using 'ecrecover'
* The resolution logic is now more complicated by the fact that there is a built-in 'challenge' or 'dispute' period every time new state is submitted.


### To make it easier

to standardize the generic state channel functionality in a way that cleanly splits the **state channel resolution logic** from the application logic.

One way to accomplish is to **model an application as a state machine**

We need a standary way to **interface with the state transitions of your application** regardless of the implementation of public API.


### How to do

Requirements:

* state transition
* authorization logic
* resolution logic

Natually, the state channel object should *use* the application logic to determine if a transaction is valid.

The StateChannel can use the App Logic as a means of determining valid transitions but handle the authorization and resolution logic itself based on information provided to it by the app logic.

### Sidechannel contract APIs

* Creating a dispute
	* one party submitting the latest signed copy of the state and optionally taking an action on the app that will logically progress the state to the next state.
* Progressing a dispute
	* One party responding to the dispute from another party by taking an action on the state that has been submitted.
* Cancelling a dispute
	* Both parties agreeing to cancel the dispute and resume off-chain normally.

#### Where are the funds stored?

generic multisignature wallet, 

and use it as the primary contract that makes commitments for distributing those funds based on the outcome of the state channel.

### Summary

* state transition
	* with the help of app logic
* authorization
	* ecrecover
* resolution
	* on-chain multi-signature wallet
	* dispute handling

