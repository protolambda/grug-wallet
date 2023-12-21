package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type RawID string

// 32 hex encoded bytes, with 0x prefix, with quotes, is the maximum length
const maxIDLength = 32*2 + 2 + 2

var _ json.Marshaler = RawID("")
var _ json.Unmarshaler = (*RawID)(nil)

// isValid checks if the JSON RPC message identifier is valid.
// if not null, a string, or an integer number, then invalid.
// I.e. leading whitespace is not valid, true/false are not valid, floats are not valid, maps and arrays are not valid.
func (id RawID) isValid() bool {
	if len(id) == 0 || len(id) > maxIDLength {
		return false
	}
	if id == "null" {
		return true
	}
	if id[0] == '"' { // any string without whitespace
		return id[len(id)-1] == '"' && json.Valid([]byte(id))
	}
	// any integer number, but no "1e9" etc.
	for _, c := range id {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (id RawID) MarshalJSON() ([]byte, error) {
	if id.isValid() {
		return nil, fmt.Errorf("invalid ID: %x", []byte(id))
	}
	return []byte(id), nil
}

func (id *RawID) UnmarshalJSON(data []byte) error {
	if id == nil {
		return errors.New("cannot unmarshal into nil RawID")
	}
	*id = RawID(data)
	if !id.isValid() {
		return fmt.Errorf("invalid ID: %x", data)
	}
	return nil
}

func (id RawID) String() string {
	return string(id)
}

type Request struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params,omitempty"`
}

type ErrorObject struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Response struct {
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *ErrorObject     `json:"error,omitempty"`
}

// V2 is a zero-cost constant type, for encoding/decoding JSON-RPC messages:
// it validates the JSON-RPC version, without allocating it as Go string in every message.
type V2 struct{}

func (V2) MarshalJSON() ([]byte, error) {
	return []byte("2.0"), nil
}

func (*V2) UnmarshalJSON(data []byte) error {
	if !bytes.Equal(data, []byte("2.0")) {
		return fmt.Errorf("invalid JSON RPC version: %x", data)
	}
	return nil
}

type Message struct {
	*Request
	*Response
	ID      RawID `json:"id"`
	JSONRPC V2    `json:"jsonrpc"`
}

func (m *Message) RespondSuccess(data any) (*Message, error) {
	if m.Response != nil {
		return nil, fmt.Errorf("cannot respond to a response: %s", m.ID)
	}
	x, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode response: %w", err)
	}
	result := json.RawMessage(x)
	return &Message{
		Request: nil,
		Response: &Response{
			Result: &result,
			Error:  nil,
		},
		ID:      m.ID,
		JSONRPC: V2{},
	}, nil
}

func ConstErrorObj(c ErrorConst) *ErrorObject {
	return &ErrorObject{
		Code:    c.Code(),
		Message: c.Message(),
		Data:    nil,
	}
}

func AnnotatedErrorObj(c ErrorConst, err error) *ErrorObject {
	return &ErrorObject{
		Code:    c.Code(),
		Message: c.Message() + ": " + err.Error(),
		Data:    nil,
	}
}

func (m *Message) Respond(data any) *Message {
	resp, err := m.RespondSuccess(data)
	if err != nil {
		return &Message{
			Request: nil,
			Response: &Response{
				Result: nil,
				Error:  AnnotatedErrorObj(InternalError, err),
			},
			ID:      m.ID,
			JSONRPC: V2{},
		}
	}
	return resp
}

func (m *Message) RespondErr(errObj *ErrorObject) *Message {
	return &Message{
		Request: nil,
		Response: &Response{
			Result: nil,
			Error:  errObj,
		},
		ID:      m.ID,
		JSONRPC: V2{},
	}
}

type RawConnection interface {
	Write(msg *Message) error
	Read(dest *Message) error
	Close() error
}

type Channel struct {
	Inward  <-chan *Message
	Outward chan<- *Message
}
