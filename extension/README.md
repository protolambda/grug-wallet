# extension

This is the extension that plugs into the browser,
to make `window.ethereum` available.

Quick start, what does this extension do:
- Go server lists provider definitions
- A provider definition has the EIP-6963 `info`: name / icon data for the RPC.
- A provider definition has an `endpoint`, an RPC endpoint to connect to as wallet.
- The Go server can provide many such configurable RPC endpoints on demand.
- A provider definition can be instantiated into an EIP-1193 provider
- The EIP-1193 provider is what exposes the RPC in the dapp page.

The goal is to make the most light-weight extension that can still:
- provide human-readable dynamically generated provider definitions to the browser
- give the user maximum choice over which provider they use
- maintain the actual provider RPC clients in the background, to not kill/restart every page change.

## Test runs

1. [Install web-ext](https://extensionworkshop.com/documentation/develop/getting-started-with-web-ext/)
2. In `./extension` run `web-ext run`, this will create a live-reloading firefox web-extension dev environment
3. Navigate to `about:debugging#/runtime/this-firefox` in firefox to open the extension debug/inspect view.
4. Open any regular web-page that does not have special rules against extension web-page content modification
5. Click the extension thing, pin the icon, and try add a server endpoint in the UI.

TODO: the providers retrieval from server endpoint has not been implemented yet.
A hardcoded default providers response is provided instead.

## Sandboxes

The web-browser is complicated. Many sandboxes, although for good browser-security reasons.

There are 4 different sandboxes involved:
- the in-page scripts, where the dapp lives
- the content-script, where we insert the provider with
- the background-script, where we maintain our RPC websockets
- the popup-script, where we configure the provider-sources

To keep provider management sane, we use the EIP-6963 provider-detail UUID
as the canonical way to identify providers across the stack.
The provider info is in the extension local storage, accessible across the 3 extension scripts.

The content-script inserts into the page by using `window.wrappedJSObject`:
the content-script has to explicitly inject the functions into the lower-security-privilege in-page script.
This works well on firefox, and prevents all the inefficiency that comes with that alternative
of embedding an in-page script that talks to the content-script over yet another stream.

The content-script has a single type of EIP-1193 provider: the `BackgroundProvider`.
This provider simply streams the back/forth messages to the background-script.

The background-script automatically creates Websocket providers to serve each provider.
In the future we may add more types of providers.
This websocket provider then communicates with the Go RPC server.

The popup script is the configuration entry-point: the user enters the server endpoint,
and the background script will be notified to load the providers from it.
The new providers then get communicated back to the popup, for a UI update.

The popup allows providers to be disconnected, re-advertised to the dapp,
and forced into the `window.ethereum` place.


## manifest

The manifest defines the extension.

[`manifest.json` permissions docs](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/manifest.json/permissions)

permissions are a work in progress. Running RPCs in the background requires some permissions.
Especially if you don't want to click the extension icon every time the page navigates somewhere.

Perhaps with `optional_permissions` and deferred content-scripts we can reduce the required permission at runtime.
This adds even more complexity however.
See:
- https://extensionworkshop.com/documentation/develop/request-the-right-permissions/#request_permissions_at_runtime
- https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/scripting/registerContentScripts


## EIPs

### 1193

A JS interface for a JSON-RPC provider. 

EIP-1193 TLDR:
- `Provider.request(args) -> promise<unknown>` for RPC calls
  - args type: `RequestArguments`: `{method: string, params?: unknown[] | object}`
  - result type: `{code: number, data?: unknown, message?: string}`
  - reject with `ProviderRpcError` standard data
    - if disconnected: error=4900
    - if connected to wrong chain: error=4901
    - if rpc method is unknown: error=4200
    - if user rejected: 4001
    - if unauthorized: 4100
- Events:
  - `Provider.on(eventName: string, listener: function)`
    - `message`
      - `method: eth_subscribe` results in `eth_subscription` results
        - result: `{type: const string, data: {subscription: string, result: unknown}}`
    - `connect` result: `{chainId: string}`
    - `disconnect` result: `{error: {code: number, data?: unknown}}`
    - `chainChanged` result: `{chainId: string}`
    - `accountsChanged`: result `{accounts: string[]}`
  - `Provider.removeListener()`

And commonly wallets race to make their provider available as `window.ethereum`.
And an `ethereum#initialized` window event is fired when the global is ready.

### EIP-6963

This EIP allows the wallet to announce its provider(s).

Summary of important types:
```
interface EIP6963ProviderInfo {
  uuid: string; // UUIDv4 to distinguish different EIP-1193 sessions
  name: string; // Grug Wallet
  icon: string; // URI-encoded image
  rdns: string; // reverse dns, e.g. com.protolambda.grugwallet
}

interface EIP6963ProviderDetail {
  info: EIP6963ProviderInfo;
  provider: EIP1193Provider;
}

interface EIP6963AnnounceProviderEvent extends CustomEvent { // dispatched by wallet
  type: "eip6963:announceProvider";
  detail: EIP6963ProviderDetail;
}
interface EIP6963RequestProviderEvent extends Event { // dispatched by dapp
  type: "eip6963:requestProvider";
}
```
The two-way request/response dispatch ensures that regardless of dapp/wallet init order,
the providers are detected.


### Unofficial

"Connectors" are a pattern that allow dapps / dapp-frameworks to get wallet providers into their dapp.
The connectors commonly send a `wallet_addEthereumChain` RPC request on initialization.

`window.etherum.isMetaMask` is set on the "real" metamask. Can add option to pretend to be that.

`window.CoinbaseWalletProvider` is like `window.ethereum` but custom.
`window.CoinbaseWalletSDK` appears to be similar like the above.
`window.ethereum.isCoinbaseBrowser`: hint for EIP-3085 (`wallet_addEthereumChain`) network switcher menus.
`window.ethereum.isCoinbaseWallet`: hint for other extra functionality.

Common RPC methods:
- `eth_accounts`
- `net_version`
- `eth_chainId`
- `personal_sign`
- `eth_sign`
- `eth_signTypedData`
- `eth_signTypedData_v4`
- `eth_sendTransaction`
- `eth_blockNumber`
- `eth_getBalance`
- `eth_getCode`
- `eth_getTransactionCount`
- `eth_getStorageAt`
- `eth_getBlockByNumber`
- `eth_getBlockByHash`
- `eth_getTransactionByHash`
- `eth_getTransactionReceipt`
- `eth_estimateGas`
- `eth_call`
- `eth_getLogs`
- `eth_gasPrice`
- `wallet_getPermissions`
- `wallet_requestPermissions`
- `wallet_watchAsset`
- `safe_setSettings`
- `metamask_accountsChanged`
- `metamask_unlockStateChanged`
- `metamask_chainChanged`
- `metamask_getProviderState`

These RPC methods should be supported by the Go server, not the browser extension.

