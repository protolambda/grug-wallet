package server

import (
	"errors"
	"github.com/protolambda/grug-wallet/server/jsonrpc"
	"golang.org/x/net/websocket"
	"io"
	"sync/atomic"
)

// TODO register like http.Handle("/websocket", websocket.Handler(WebSocketHandler))
// TODO also add some validation / auth middleware
func WebSocketHandler(ws *websocket.Conn) {
	ws.MaxPayloadBytes = 32 << 20

	// note: Receive(ws,msg) calls an internal frame handler of the provided ws,
	// which should handle ping/pong/close frames that may be interleaved with the actual data.

	// TODO create WS
}

// WS represents a websocket JSON RPC connection. It can be used for both client-side and server-side
type WS struct {
	ws        *websocket.Conn
	connected atomic.Bool
}

// Write to the RPC, safe for concurrent use
func (w *WS) Write(msg *jsonrpc.Message) error {
	if err := websocket.JSON.Send(w.ws, msg); err != nil {
		if errors.Is(err, io.EOF) {
			w.connected.Store(false)
		}
		return err
	}
}

// Read from the RPC, safe for concurrent use
func (w *WS) Read(dest *jsonrpc.Message) error {
	if err := websocket.JSON.Receive(w.ws, dest); err != nil {
		if errors.Is(err, io.EOF) {
			w.connected.Store(false)
		}
		return err
	}
}

func (w *WS) Close() error {
	if err := w.ws.Close(); err != nil {
		return err
	}
	w.connected.Store(false)
	return nil
}

var _ jsonrpc.RawConnection = (*WS)(nil)
