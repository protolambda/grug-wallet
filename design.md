# Grug wallet design

## Background

There are so many "web3" wallets out there!
Extensions! Mobile apps! Wow.

But they all feel the same:
- Log in to account(s)
- Big central token swap button, convenient...
- Insert account into dapp
- Show balances
- Show last few trades
- Sign transactions

As a long-term user you eventually get tired of the above: it does not work in your favor.
It's a mix of hot-keys, swap bloatware, account management, and privacy issues if you look closer at it.

### Issues

Most wallets seem to get so many things "wrong" (from my perspective):
- Bad defaults:
  - Dapp asks for account, you "accept", default global account selected: missing dapp encapsulation.
  - Dapp asks for network, you "accept", RPC implied by wallet: missing RPC customization.
- Privacy:
  - Wallet vendor tracks the sites you are connected to
  - Wallet vendor tracks the trades you make
  - Wallet vendor tracks the accounts you own
  - Sites track that you are a web3 user
  - "read-only accounts" still a rare feature
  - No proper domain separation between wallets
  - Wallet-connect is cool, but centralized wallet relays are not
  - Why does the wallet rely on centralized indexing?
- Orderflow:
  - Wallet suggests the RPC you use, does not let you maintain multiple RPCs
  - Wallet suggests the swap you use, does not let you replace / remove the feature
  - Wallet uses the same RPC for reads (calls, gas estimates) and writes (txs),
    exposing probabilistic intents (MEV!) and making intentional routing (MEV protect!) hard.
- Login
  - Why do I need a seed-phrase for login when I connect a completely different hardware wallet?
  - Why is multitasking so difficult?
    - Can I not log in to two different accounts on two different sites?
    - Can I not use more than one RPC per network?
    - Can I not use more than one network per dapp?
  - You know about HD paths? Why can I not import accounts with custom HD paths?
- Storage
  - Local keystores impossible to find
  - Account list impossible to find
  - Why can I not import/export my interactions?
- UX
  - Why can I not persist and privately share address labels, without leaking them to a centralized service?
  - Why can I not label *other* addresses outside the wallet.
    - As a degen, how can I track the OPP?
    - As a user, how can can I track my accounts?
  - Coloring / noting addresses: apparently online poker UX is ahead of web3 UX.
    Auto-generated "blockie" address "profile-pics" were a start,
    but why not let the user customize something simple like a color, or add a note? 
  - Why does the wallet touch the fee / nonce, after I changed it?
  - Why is overriding the nonce so difficult?
  - Why can I not replay existing transactions?
  - Many extensions have a pop-up window only, and hide and reset their menu navigation as soon as it loses focus.
  - Impossible to copy-paste a full address, tx-hash or anything without opening etherscan first.
- Data:
  - There exist more assets than ERC20 and ERC721
  - There are many fake / spam assets in ERC20 / ERC721 forms, I do not need to see/interact with those in my wallet
  - It should be possible to monitor more arbitrary things:
    - Things that I am interested in should be expressed like dotfiles: shareable, customizable, extensible
    - Collecting the data should not require centralized indexers
    - Data should not spam an RPC, but rather update when the user wants it to update, and clarify its origin.
- Attacks:
  - Wallet pulls in so many NPM dependencies, a scary supply-chain really
  - Wallet exposes itself so deeply in the browser
  - Difficult/impossible to track your interactions with it.
    - The wallet analytics share more about you with their developer than they share with you
  - Difficult/impossible to audit the RPC requests by a page
  - Difficult/impossible to hold an RPC request before ACK.
  - Difficult/impossible to queue up a tx. Signing requests in your face!
    Easy to sign the wrong stuff with UI exhaustion / spam.
  - Difficult to track past signature requests (unsigned transactions) and transactions.
  - Wallet fully trusts the RPC responses,
    does not verify transactions are correctly represented like they are included on-chain.
  - Wallet fully trusts the RPC is following the canonical chain,
    does not verify you are on the canonical chain.
- Bloat, wallets trying to be platforms:
  - Integrated NFT gallery, while there are 10+ dapps that try to provide a gallery.
  - Integrated chat, because everything must be social, even my finances /s.
  - Integrated swapping, this is really here because of retail orderflow.
  - Integrated app store, full of ads for dapps that struggle to find PMF.
- Hardware wallets:
  - Browser API to interface with hardware devices is a joke
  - Debugging hardware wallet failures is a nightmare. Broken HID interface? Udev rules wrong? Let's not tell the user.

I shared a lot of these issues in [a tweet](https://twitter.com/protolambda/status/1728823287646499154).
It made some wallet vendors unhappy. Some rightfully so, some are actually building good products (hardware wallets),
or focused on onboarding non-technical users. Others are just copy-cats of metamask. I need something different than that.

Everyone wants to build the crypto everything app. Become a platform. I get it, every crypto dev does.
But I want that platform at a safe distance from my keys, to respect my privacy, to not lock me in, 
and not assume my orderflow. Crypto should be a tool to serve, not to capture, the users.

### Embedded wallets

This is an interesting design route, especially for its user-friendliness.

In this model, the user only has to provide a signer. The rest of the wallet, is embedded into the dapp.
The user provides an account signer specific to that dapp.

This signer is minimal, and can be native to the browser (passkeys! mobile enclaves!).
The user does not need to install a separate wallet to start using the dapp.
The user should get hardware-level security by default when onboarding to the dapp.

If the user trusts the dapp frontend, and the signers are separate per dapp, then their assets are generally safe:
only if the dapp frontend is malicious then they lose assets.

Combine this with immutable frontends, and the user is in a great place:
- they don't put up all their assets into one account for grabs
- they can go back to the frontend any time
- the frontend can serve all dapp-specific inventory logic (positions, assets, past actions, etc.)

Also, if you don't recognize the design parallel already:
- Concurrency is done through communication, not state sharing.
- Interoperability is done through transactions/intents, not account sharing.

Where embedded wallets fall short:
- The EVM cannot validate common types of signatures, not all keys work (e.g. no yubikey transactions yet).
- If the frontend is gone, most users won't know what to do to use their assets.
- Reduces bloat in the wallet, but adds it to the dapp:
  the dapp must do a really good job to make it transparent and trustless.
- Onboarding also requires funding of fees: funding an embedded wallet means coupling more external things to the dapp,
  to try and get the user to pay for their account.
  Onboarding through a credit-card is not decentralized, so you need sponsored transactions and/or external wallets.
- Without account-abstraction, the key becomes irreplaceable.

### Account abstraction

"Account-abstraction" is challenging, and we are seeing two different variants:
Separating execution from accounts:
- offchain, EIP-4337: with user-ops, bundlers and paymasters: by defining the intents as user-ops,
  to be bundled into a tx by any account, and paid for by user to bundler through a paymaster, akin to a bounty.
- onchain, EIP-3074: with `AUTH`/`AUTHCALL` EVM opcodes: by delegating the EOA identity to invoker contracts.

These two approaches can complement each other too.

**[EIP-4337](https://www.erc4337.io/docs)** is not really native "abstraction",
as it pushes everything to the app/infra-layer, adding lots of wrapper inefficiency.
But at least it keeps the protocol simple: no EVM changes are required.

You may have seen twitter-threads how expensive the 4337 transactions
get when the L2 transaction-fees are basically dominated by data costs,
while 4337 adds significant metadata and signature data to each transaction.

The 4337 abstraction does make sense though:
"users don't pay AWS fees for using web2 apps, so paying gas fees for using dapps could feel foreign and wrong to them".

**[EIP-3074](https://eips.ethereum.org/EIPS/eip-3074)** is "native", but comes with the risk of delegating *too much*.
You have to trust the contract you are delegating too, or it can maliciously act on your behalf.

This EIP requires breaking changes to the EVM, and has not been included in mainnet ethereum yet.

It comes with interesting new abilities too though:
- Parallel nonce values
- Commit to multiple calls with a single signature

Account abstraction can offer great things:
- Sponsored transactions are a great onboarding experience
- Decoupling account keys from account identity is a great security tool (key rotation!)

Ok, and now back to real-world:
- EIP-4337 adds overhead (may get better on L2 when L1 DA gets cheaper), separates user more from app.
- EIP-3074 is not live, although it has been around for very long, so much so that it is becoming a meme.

The effect on browser/mobile wallets? 
Like embedded wallets, it allows the dapp to bring more of the wallet features, now including transaction sponsoring.
But with abstraction comes complexity, something to stay alert of.

It is magical to send a friend a link, and the link onboards them right on to the dapp,
but I feel more transparency and explicit key-ownership can help the power-user:
if it's not about a simple coffee-size payment I want to know what's happening.

### Legacy and universal accounts

Not all "dapps" are the same. New social / integrated experiences can benefit from an embedded wallet
and better onboarding flow, giving the user a quick dapp-specific identity.

But many of the dapps that have found PMF with the power-users are not like that.
The web3 masterminds have taken "sovereignty of your data" and turned it into "use the same account everywhere".

Reusing an account is quite common in ethereum, even though in all honesty it is bad practice.

And many legends/OGs have popularized this: ENS, Status and Web3js have designed around this. 
In my opinion, it has a place, but fits better in the social sphere, not all of crypto.

"universal accounts" can kind of resolve this:
- Dapps can all have their own hotkeys
- you can use the hot-key in the embedded wallet
- The hot-key has limited access to a smart-contract account
- The smart-contract account carries the legacy of having a singular named social account.

I proudly have an ENS, and like that crypto has a social ecosystem,
but my maybe unpopular opinion is that not everything needs to be coupled back to a "universal account".
It's bad for privacy, opsec, and managing all those inter-account relationships
onchain is more trouble than it is worth to me. 

### Multisig wallets

Multisigs are more like embedded wallets again: they require N signers,
and then present the wallet-interface as a dapp frontend.

The individual signers often are backed by regular browser/mobile wallets,
which then couple back to the actual signers (hardware wallets even).

Multisigs can be a great intermediate wallet layer to improve signer security,
but do not fix many of the wallet UX issues: multi-accounts, RPC, attack-surface issues are still real.


### Why not use a CLI wallet, you are technical, right?

The "web3" term has become a buzz term, but does get one thing right: the web is the universal interface.

Through this interface we can:
- Change application logic, while providing the same familiar UI.
- Add new functionality, while making it discoverable in the same UI.
- Add app-layer user-friendly interfaces, for highly technical state changes.

All of the above are not possible with a CLI wallet.

Instead, you end up paying a cost for premature automation, and manually entering inputs into contracts is error-prone.
Even worse, often it requires reverse-engineering, as the smart-contract interface is more like a plan B,
and the user is expected to use the UI in the happy path.


## Becoming a Grug wallet developer

I do not really know how to fix web3 wallets, I am not a "wallet" developer.
Onboarding non-technical users is really hard. It is not the problem I will solve here.

But I do want to have a working wallet, or I miss out on the exciting and degenerate world that is onchain.
Existing wallets have been a pain to use, they do not do what I want them to do.

So here I am, figuring out how to become a [Grug](https://grugbrain.dev/) wallet developer:
- I know nothing.
- I will say "no" to many existing in-wallet features and ideas (social, swaps, anything Typescript).
- I will build what works for me.
- I will use tools that work for me.
- I will break apart the "wallet" concept as much as I like: factor your code.

Repeat after me: "complexity is bad".

So what ideas do I have, to remove complexity?
- Minimal/no Javascript
- Not in the browser/mobile, but isolated
- Focus on configurable essentials (keys, RPCs, networks, transactions)
- No integrated bloatware

The key here is that the *user brings their own configs* (dotfiles, datasources).
- The wallet just needs to understand sources and config types.
- It does not need to implement asset indexing, NFT galleries, or anything like that. 


### Design

Just because the "web3" is local-first, does not mean we cannot use backend tech.
Rather, it's the opposite: if there is anything the user should self-host, it is the backend.

Running a full node is maybe not for every user, but running a light-client is.
This solves many of the offchain privacy and security issues, by not relying on centralized trusted defaults.

Now, the next step is to run the wallet like this: decoupled, free from the browser.

#### HTMX

And if you like Grug, and dislike javascript, you will like [HTMX](https://twitter.com/htmx_org).
Actually, the [HTMX docs are here](https://htmx.org/) but I figured their twitter is a vibe you may also like.

With a local server, you can open localhost in your browser, and you will have a wallet interface that just works.
- No pop-up issues
- No focus-loss issues
- No browser-local storage issues
- No browser hardware-wallet interface issues

And, thanks to HTMX: no javascript, no NPM!

And grug does not care: serving the interface on localhost is a feature, not a bug.
Wrapping it in something like electron is ridiculous.

And HTMX works with any backend: it is a "high-power tool", not a "framework".

Personally, this Grug prefers Go as programing language:
it is familiar, and if it is good enough for go-ethereum, then it is good enough for Grug.
Besides, Go has proven to be capable of running the kind of configuration-first servers that I like,
like Grafana and Prometheus.

#### Extension

Now, we do want the browser-extension still: without this,
we need to rely on something like walletconnect to couple to dapps,
and can't do the RPC integration like we want to do.

So a little javascript will be necessary, to define the firefox/chrome/etc. browser plugin definitions.
There should not be a need to even use NPM here.

But this can be the bare minimum: just a tool to act like metamask, configure the local endpoint,
but then pass on all the interactions to the server on localhost.

#### Mobile

And maybe the same can be done for mobile: a bare minimum app,
that provides the phone with the web3 "intent" handling of the PWA
(progressive web app, the hack where all crypto dapps run within).

It may even mean you can use a mobile dapp, while managing your account on browser, and signing with a hardware wallet!
No need for fancy bluetooth wallets or mobile hotkeys.

#### Account separation

The server can handle the calls based on the caller origin / headers, to differentiate between concurrent dapps.
Like embedded wallets, separate accounts can be generated per dapp.

#### Configuration-first

The last step is all about configuration:
accounts, RPCs, networks, token-lists, contract-sources, balance lookups, labels, etc.
are all just different dotfile-like data-sources!

- Accounts should each be a metadata file,
  with details where the key-material can be found (hotkey, keystore, HD path, remote, hardware, etc.).
- RPCs should each be a metadata file, with details like the endpoint, rate-limit, chain ID.
  - Some RPCs may be defined as combinations of other RPCs:
    e.g. all archive requests to remote endpoint X, all regular requests to local endpoint Y,
    and all tx requests to MEV-protect endpoint Z.
- Networks should each be a script or metadata file,
  with details like network name, currency name, block-explorer options.
  - The "add network" flow would not be necessary if we can just load network metadata lists from git repositories.
- Token-lists should each be a script, that can then lookup token details.
  There exist different token-list repositories / APIs, and a script can provide an interface to that data.
- Contract ABI are also just scripts:
  - there may be an etherscan-API powered source
  - there may be a local DB
  - there may be various local reposities with build artifacts
  - there may just be a JSON file
- Account status is also just scripts:
  - ERC20 balance lookups may be defined by the user as `cast` commands.
  - More advanced position monitoring may be defined by the user as foundry scripts.
- Addresses, colors, labels, notes should be configurable

And for each of these, there may be overlap: as a wallet we can flag duplicates, and deduplicate 1:1 copies.
But allowing multiple RPCs per network is a feature, not a bug.

And ideally, dotfiles can be sourced from multiple directories, that all have the same structure:
- a source may be private, or public (git repositories!)
- sources may have overlapping and different data

If we get the data-source pattern right, then the
*wallet does not need to know all that much about the application layer*.

All the wallet has to do is load all configurations, merge them, and combine them
to provide a usable wallet experience.

#### Prototype

To sum up, for a prototype, we will need:
- A minimal browser extension definition, that:
  - Has a configurable server endpoint
  - Acts like `window.ethereum` (or other "modern" wallet coupling approach), to pass on any interactions to server.
- A Go server that:
  - Loads data-sources
    - Executes scripts, parses non-scripts
  - Provides an API for the browser-extension to function
  - Provides an HTMX GUI for browsing the datasources and handling the tx-approval flow

This is not too bad: we can start with the bare minimum browser-extension / server separation,
define some minimal data-sources for accounts and RPCs,
and start building an HTMX flow to handle transaction signing.

