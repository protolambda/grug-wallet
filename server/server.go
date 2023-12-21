package server

import (
	"context"
	"github.com/protolambda/grug-wallet/server/jsonrpc"
)

type UUIDv4 string

type EIP6963ProviderInfo struct {
	UUID UUIDv4 `json:"uuid"` // UUIDv4 to distinguish different EIP-1193 sessions
	Name string `json:"name"` // Grug Wallet
	Icon string `json:"icon"` // URI-encoded image
	RDNS string `json:"rdns"` // reverse dns, e.g. com.protolambda.grugwallet
}

type EIP6963ProviderEndpoint struct {
	Info     EIP6963ProviderInfo `json:"info"`
	Endpoint string              `json:"endpoint"` // 'ws://127.0.0.1:8646' TODO
}

type EIP6963Provider struct {

	// TODO websocket server
}

type EIP6963Providers struct {
	Providers map[UUIDv4]*EIP6963Provider
}

// TODO EIP-6963 providers pool

// TODO collection of active RPC connections
type OutgoingConnection struct {
	Conn jsonrpc.RawConnection
}

// TODO turn websocket into RawRPCConnection

// TODO turn RawRPCConnection into RPCChannel

type IncomingConnection struct {
	Conn RPCChannel
}

// TODO configuration resource pool

type IncomingClientID string

// TODO backend

type Backend interface {
	RegisterRPC(ctx context.Context, from IncomingClientID, prov UUIDv4)
}

type incomingReq struct {
	conn *IncomingConnection
}

type Server struct {
	// connections we maintain to external RPC servers
	Out map[*OutgoingConnection]struct{}
	// connections we have to handle RPC clients
	In map[*IncomingConnection]struct{}

	incoming chan incomingReq
}

func (s *Server) mainLoop() {
	for {
		select {
		case req := <-s.incoming:
			s.In[req.conn] = struct{}{}

		}
	}
}

// websocket upgrade HTTP route
//  -> parse UUIDv4 from URL, maybe do some validation
//  -> upgrade to websocket, register websocket globally
//  -> setup websocket RPC read/write loop
//  -> register RPC with backend

// backend:
//  -> setup RPC handler: read/write from RPC currently associated with the prov, translate message IDs
//    -> redirect some requests to the signer queue
//    -> buffer some requests for ACK

// TODO HTTP route for EIP-6963 providers list
// TODO HTTP routes for each EIP-6963 RPC
// TODO HTTP routes for icon assets
// TODO HTTP route for HTMX frontend

// TODO CLI frontend maybe
