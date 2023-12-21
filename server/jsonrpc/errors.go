package jsonrpc

import "fmt"

type Error interface {
	Code() int64
	Message() string
}

type ErrorConst int64

const (
	// Invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	ParseErr ErrorConst = -32700
	// The JSON sent is not a valid Request object.
	InvalidRequest ErrorConst = -32600
	// The method does not exist / is not available.
	MethodNotFound ErrorConst = -32601
	// Invalid method parameter(s).
	InvalidParams ErrorConst = -32602
	// Internal JSON-RPC error.
	InternalError ErrorConst = -32603
	// (EIP-1474) Missing or invalid parameters
	InvalidInput ErrorConst = -32000
	// (EIP-1474) Requested resource not found
	ResourceNotFound ErrorConst = -32001
	// (EIP-1474) Requested resource not available
	ResourceUnavailable ErrorConst = -32002
	// (EIP-1474) Transaction creation failed
	TransactionRejected ErrorConst = -32003
	// (EIP-1474) Method is not implemented
	MethodNotSupported ErrorConst = -32004
	// (EIP-1474) Request exceeds defined limit
	LimitExceeded ErrorConst = -32005
	// (EIP-1474) Version of JSON-RPC protocol is not supported
	JSONRPCVersionNotSupported ErrorConst = -32006
	// (EIP-1193) The user rejected the request.
	UserRejectedRequest ErrorConst = 4001
	// (EIP-1193) The requested method and/or account has not been authorized by the user.
	Unauthorized ErrorConst = 4100
	// (EIP-1193) The Provider does not support the requested method.
	UnsupportedMethod ErrorConst = 4200
	// (EIP-1193) The Provider is disconnected from all chains.
	Disconnected ErrorConst = 4900
	// (EIP-1193) The Provider is not connected to the requested chain.
	ChainDisconnected ErrorConst = 4901
)

func (c ErrorConst) Code() int64 {
	return int64(c)
}

func (c ErrorConst) Message() string {
	switch c {
	case ParseErr:
		return "Parse error"
	case InvalidRequest:
		return "Invalid Request"
	case MethodNotFound:
		return "Method not found"
	case InvalidParams:
		return "Invalid params"
	case InternalError:
		return "Internal error"
	case InvalidInput:
		return "Invalid input"
	case ResourceNotFound:
		return "Resource not found"
	case ResourceUnavailable:
		return "Resource unavailable"
	case TransactionRejected:
		return "Transaction rejected"
	case MethodNotSupported:
		return "Method not supported"
	case LimitExceeded:
		return "Limit exceeded"
	case JSONRPCVersionNotSupported:
		return "JSON-RPC version not supported"
	case UserRejectedRequest:
		return "User Rejected Request"
	case Unauthorized:
		return "Unauthorized"
	case UnsupportedMethod:
		return "Unsupported Method"
	case Disconnected:
		return "Disconnected"
	case ChainDisconnected:
		return "Chain Disconnected"
	default:
		return fmt.Sprintf("Non-standard error-code %d", c.Code())
	}
}

// IsServerError identifies server errors, per standard JSON-RPC 2.0 error code scheme.
// Reserved for implementation-defined server-errors.
func (c ErrorConst) IsServerError() bool {
	return c < -32000 && c > -32099
}
