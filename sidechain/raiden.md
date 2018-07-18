## Raiden

```
CreateNetwork (address token_address) public

event NetworkCreated(token_address, token_network_address)
```


```
OpenChannel (participant1, participant2, settle_timeout) public returns (channel_identifier)

event ChannelOpened(identifier)
```


```
SetDeposit(participant, total_deposit, partner) public

event ChannelNewDeposited(channel_identifier, participant, deposit)
```

```
Withdraw(participant, withdraw, partner, []signatures)

event ChannelWithdraw
```

```
CloseChannel(partner, balance_hash, nonce, additional_hash, signature)

event ChannelClosed(channel_identifier, []participants)
```

```
SettleChannel([]participant, []amount, []root) public

event channelSettled(channel_identifier, []participant_amount)
```