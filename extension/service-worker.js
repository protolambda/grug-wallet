// This is the service worker script, which executes in its own context
// when the extension is installed or refreshed (or when you access its console).
// It would correspond to the background script in chrome extensions v2.

console.log("Hello world grug wallet")

importScripts('service-worker-utils.js')

// Note all import paths are relative to the file `manifest.json`.
