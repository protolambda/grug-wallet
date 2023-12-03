
browser.runtime.onMessage.addListener(async (message) => {
    console.log("background-script received message", message);
    if (message.command === "loadProviders") {
        await loadProviders(message.server);
        await browser.runtime.sendMessage({
            command: "refreshProviders",
        });
    }
    if (message.command === "disconnectProvider") {
        const prov = openProviders[message.providerUUID];
        if(prov) {
            prov.disconnect();
        }
        // TODO: fire UI event
    }
});

// loads the provider data from the server, and puts it all in session storage for later use.
async function loadProviders(server) {
    console.log("fetching providers from server", server);

    // TODO: actually fetch real provider definitions from the server

    const providerName = "todo";

    const image = 'data:image/svg+xml,<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="32px" height="32px" viewBox="0 0 32 32"><circle fill="yellow" cx="16" cy="16" r="12"/></svg>';

    const providerUUID = 'baa2ce38-be93-4ec4-8b5b-1ccb5b17a43c';

    const info = {
        uuid: providerUUID, // UUIDv4 to distinguish different EIP-1193 sessions
        name: "Grug wallet "+providerName, // human readable display name
        icon: image, // URI-encoded image
        rdns: "com.protolambda.grugwallet", // reverse dns, to identify wallet type
    };

    const providerEndpoint = {
        info: info,
        endpoint: 'ws://127.0.0.1:8646', // TODO
    }

    await browser.storage.local.set({
        ['web3-'+providerUUID]: providerEndpoint,
    })
}

class JSONRPCWebSocket {
    constructor(url) {
        console.log("initializing websocket", url);
        this.url = url;
        this._listeners = [];
        this._connect();
    }

    _connect() {
        if(!this.url) {
            console.log("aborting connection attempt, no url specified");
            return;
        }
        console.log("connecting websocket client to:", this.url);
        this._socket = new WebSocket(this.url);
        const incomingMessage = (event) => {
            console.log("background->content JSON RPC message", event.data);
            this._listeners.forEach(listener => listener({rpcEvent: event.data}));
        }
        this._socket.addEventListener("message", incomingMessage.bind(this));
        const onError = (error) => {
            console.error('WebSocket Error:', error);
            this._listeners.forEach(listener => listener({errorEvent: error.message}));
        }
        this._socket.addEventListener("error", onError.bind(this));
        const onClose = () => {
            console.log('WebSocket Closed. Need to reconnect...');
            this._listeners.forEach(listener => listener({closeEvent: null}));
            if (this.url) {
                this.reconnectTimer = setInterval(() => {
                    this._connect();
                }, 10000); // reconnect after 10 seconds
            }
        }
        this._socket.addEventListener("close", onClose.bind(this));
        const onOpen = () => {
            if(this.reconnectTimer) {
                clearInterval(this.reconnectTimer);
                this.reconnectTimer = undefined;
            }
            console.log('WebSocket Opened.');
            this._listeners.forEach(listener => listener({connectEvent: null}));
        }
        this._socket.addEventListener("open", onOpen.bind(this));
    }

    attachPort(port) {
        console.log("initializing port listener", port.name);
        const outgoingMessage = (rawJSON) => {
            console.log("content->background->out JSON RPC message", rawJSON)
            this._socket.send(rawJSON);
        }
        port.onMessage.addListener(outgoingMessage.bind(this));
        this._listeners.append(port.postMessage);
        console.log("initialized port listener", port.name);
    }

    async reconnect() {
        if (this._socket) {
            this._socket.close();
        }
        this._connect();
    }

    async disconnect() {
        console.log("disconnecting");
        this.url = undefined; // prevents it from reconnecting
        this._socket.close();
    }
}

const openProviders = {};

// Long-lived communication between content-script and background-script
// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Content_scripts#connection-based_messaging

async function connected(port) {
    // p: runtime.Port
    // https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/runtime/Port
    if(!port.name) {
        return
    }
    if(!port.name.startsWith('web3-')) {
        console.log("ignoring unrecognized connection request", port.name);
        return
    }
    const providerUUID = port.name.substring(5); // strip 'web3-'
    const providerEndpoint = (await browser.storage.local.get(port.name))[port.name];
    console.log("prov end", providerEndpoint);
    console.log("initializing provider", providerEndpoint.info, providerEndpoint.endpoint);
    // TODO: HTTP support
    let rpc = openProviders[providerUUID];
    if (rpc) {
        rpc.attachPort(port);
        console.log("attaching to existing provider", providerUUID);
        return;
    }
    // remember the provider, so we can close it with the popup UI when we want to.
    rpc = new JSONRPCWebSocket(providerEndpoint.endpoint);
    openProviders[providerUUID] = rpc;
    rpc.attachPort(port);
    console.log("initialized provider");
}

browser.runtime.onConnect.addListener(connected);

