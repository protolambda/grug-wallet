# Layouts

*braindump of ideas, very experimental*

## Home

- source RPCs
  - connection status
  - endpoint
  - vendor

- composite RPCs

- provider RPCs
  - client count
  - open tasks

- wallets
  - accounts

- account search

- hardware list
  - connections status

- chains

- global transaction log
- global source RPC log
- global provider RPC log

- button to add source RPC
- button to create composite RPC
- button to create wallet (list of accounts)
- button to add chain
- button to add hardware
- button to create signer
- button to create provider


## Provider

- current chain
- RPC routes
- Clients list
- Open Accounts list
- Tasks list
- transaction log
- provider RPC log
- buttons to create new tasks
  - send tx task


### Task: sign tx

show transaction type - button to change it
show transaction attributes - edit option each

reset button to revert to original tx request

account-selector button to find address for "to"

calldata ABI decoder

gas price with fast/slow/0.5/2x button

priority fee button

clear Chain ID, red warning if unexpected chain,
or unspecified chain (legacy tx)

live updating ETH fee cost

live updating fee in dollar terms

nonce update button, based on tx-pool or chain

button to calculate estimated L1-data fee.

"simulate tx" button

"sign and queue" or "sign and publish"

### Task: access account

Simple "allow" / "deny" option.
Claerly show which account.

### Task: publish tx

Show full tx

Publish a previously signed tx.

### Tx inspector

Show full tx:
Encoded tx
Decoded tx

Full color-coded Tx-hash
Buttons to open in explorer(s) of chain

Show tx fees vs onchain fees vs tx-pool fees vs suggested fee

Trace button

Simulate button

Publish button

"Fork" button -> open signing task, with preset values that match 1:1.
Can then replay with different fees, bump nonce, etc.



### Trace inspector

-> inspect raw state changes
-> inspect token-balance changes
-> inspect CALL trace
-> inspect heatmap of opcodes
-> inspect return-data (revert reason!)


## Account details

- Color-coded check-summed address
- List of wallets it is part of
- List of hardware it is signable by
- List of notes about the account

- Latest known nonce per chain, with update button

- Buttons to open in explorer(s)

- Scripts to run on account


# TODO

- maybe support bundling of txs?
- 