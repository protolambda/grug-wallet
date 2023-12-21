package server

// TODO: implement the tx handler:
//  - hold a signing request until approved
//  - hold an outgoing tx until approved
//
// eth_sendRawTransaction
// eth_sendTransaction
//	-> eth_signTransaction and then eth_sendRawTransaction
//  -> forward
