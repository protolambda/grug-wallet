(function () {
    // only inject once
    if (window.hasRun) {
        console.log("already completed inject provider work")
        return;
    }
    window.hasRun = true;

    console.log("running content script!");

    window.providerCache = {};

    // The provider we expose to the in-page scripts,
    // coupled to the background script which maintains the actual provider connection.
    //
    // This implements the EIP 1193 provider interface:
    // https://eips.ethereum.org/EIPS/eip-1193
    class BackgroundProvider {

        // port communication format:
        // content->background:  raw json-rpc as string
        // background->content: {rpcEvent: raw json-rpc string} | {closeEvent: null} |
        //                          {errorEvent: error message} | {connectEvent: null}

        constructor(providerUUID) {
            console.log("connecting");
            const port = browser.runtime.connect({name: 'web3-' + providerUUID});
            console.log("connected");
            port.onMessage.addListener(this._handleMessage.bind(this));
            this._port = port;
            this._requests = {};
            this._listeners = {
                message: [],
                connect: [],
                disconnect: [],
                chainChanged: [],
                accountsChanged: [],
            };
            this._id = 0;
        }

        async _handleMessage(event) {
            if (event.rpcEvent) {
                const data = JSON.parse(event.rpcEvent);

                if (data.id && this._requests[data.id]) {
                    const {resolve} = this._requests[data.id];
                    resolve(data.result);
                    delete this._requests[data.id];
                } else {
                    // if not matching a request, then it's a subscription response,
                    // with fields like {method, params}.
                    this._listeners.message.forEach(listener => {
                        if (data.method === "eth_subscription") {
                            listener({
                                type: 'eth_subscription',
                                data: {subscription: data.params.subscription, result: data.params.result},
                            });
                        } else {
                            console.log("received unrecognized RPC call", data);
                        }
                    });
                }
            }
            if (event.closeEvent) {
                // TODO
                console.log("content-script received close event", event);
            }
            if (event.errorEvent) {
                // TODO
                console.log("content-script received error event", event);
            }
            if (event.connectEvent) {
                // TODO
                console.log("content-script received connect event", event);
            }
        }

        async request(args) {
            const id = this._id + 1;
            this._id = id;
            const request = {
                jsonrpc: '2.0',
                id: '' + id,
                method: args.method,
            };
            if (args.params !== undefined) {
                request.params = args.params
            }

            return new Promise((resolve, reject) => {
                this._port.postMessage(JSON.stringify(request));
                this._requests[id] = {resolve, reject};
            });
        }

        on(eventName, listener) {
            this._listeners[eventName].push(listener);
        }
    }

    // Returns a EIP6963ProviderDetail
    // See EIP 6963: https://eips.ethereum.org/EIPS/eip-6963
    async function getProviderDetail(providerUUID) {
        const providerDetail = window.providerCache[providerUUID];
        if (providerDetail) {
            return providerDetail;
        }

        // the info should be in local storage
        const providerInfo = await browser.storage.local.get('web3-'+providerUUID);
        if (providerInfo) {
            const provider = new BackgroundProvider(providerUUID);

            // We need to jump through some hoops to share it with the page-content
            // https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Sharing_objects_with_page_scripts

            // Transform the promise too
            const request = (method, params) => {
                return new window.Promise((resolve, reject) => {
                    provider.request({method: method, params: params}).then((result) => {
                        resolve(cloneInto(result, window.wrappedJSObject))
                    }).catch((err) => {
                        reject(cloneInto(err, window.wrappedJSObject))
                    });
                });
            };

            // Transform the callback
            const onEv = (eventName, listener) => {
                provider.on(eventName, (ev) => {
                    listener(cloneInto(ev, window.wrappedJSObject));
                });
            };

            const providerDetail = {
                info: providerInfo,
                provider: {
                    request: request,
                    on: onEv,
                }
            }

            console.log("sharing existing provider detail", providerInfo);
            return providerDetail;
        } else {
            throw Error("unavailable provider: " + providerUUID);
        }
    }

    // shareProvider announces a single provider to the Dapp
    async function shareProvider(providerUUID) {
        console.log("announcing wallet provider", providerUUID);
        const providerDetail = await getProviderDetail(providerUUID);

        window.providerCache[providerUUID] = providerDetail;

        const announceEvent = new CustomEvent(
            "eip6963:announceProvider",
            {detail: providerDetail},
        );

        // The Wallet dispatches an announce event which is heard by the DApp code that had run earlier
        window.wrappedJSObject.dispatchEvent(cloneInto(announceEvent, window.wrappedJSObject, {cloneFunctions: true}));
    }

    // shares all known previously shared providers
    async function shareAllProviders() {
        console.log("announcing all wallet providers");
        for (const providerUUID in window.providerRegistry) {
            await shareProvider(providerUUID);
        }
    }

    // The Wallet listens to the request events which may be
    // dispatched later and re-dispatches the `EIP6963AnnounceProviderEvent`
    window.wrappedJSObject.addEventListener("eip6963:requestProvider", shareAllProviders.bind(this));

    // The classic way of injecting a provider: override window.ethereum.
    // This is not the default, but available for the clunky dapps that don't support EIP 6963 yet.
    async function overrideGlobalProvider(providerUUID) {
        console.log("overriding global window.ethereum provider", providerUUID);
        const providerDetail = await getProviderDetail(providerUUID);
        console.log("provider detail", providerDetail);

        const ethereum = createObjectIn(window.wrappedJSObject, {defineAs: "ethereum"});
        exportFunction(providerDetail.provider.request, ethereum, {defineAs: "request"});
        exportFunction(providerDetail.provider.on, ethereum, {defineAs: "on"});

        window.wrappedJSObject.dispatchEvent(new window.wrappedJSObject.Event('ethereum#initialized'));
        console.log("window ethereum var is set!", providerUUID);
    }

    // Listen for messages from the popup menu: the user can manage providers
    browser.runtime.onMessage.addListener(async (message) => {
        console.log("content-script received message", message);
        if (message.command === "shareProvider") {
            await shareProvider(message.providerUUID);
        }
        if (message.command === "overrideGlobalProvider") {
            await overrideGlobalProvider(message.providerUUID);
        }
        // TODO: maybe an option to forget providers?
    });


    // TODO: load any existing active providers

    console.log("ready to receive ethereum provider!");
})();
