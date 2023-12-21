package server

// TODO:
// RPC module that serves the offline key interactions.
// Need to decode/encode the tx.
// Compute signing-hash, to verify what's being signed with hardware-wallet.
// Air-gap support perhaps? As a form of HW wallet?

// When implementing the hold, the ctx cancel-cause can be used for signing-request replacements

// To support:
//- eth_signTypedData_v4
//- personal_sign
//- eth_signTransaction   (backs eth_sendTransaction)

// deprecated, not sure if worth supporting:
//- eth_decrypt
//- eth_getEncryptionPublicKey
