package rpc

type ErrorCode int

const (
	// ParseError indicates the request was malformed.
	ErrorCodeParse ErrorCode = -32700
	// InvalidRequestError indicates the request was not a valid JSON-RPC 2.0 request.
	ErrorCodeInvalidRequest ErrorCode = -32600
	// MethodNotFoundError indicates the method does not exist / is not available.
	ErrorCodeMethodNotFound ErrorCode = -32601
	// InvalidParamsError indicates invalid method parameter(s).
	ErrorCodeInvalidParams ErrorCode = -32602
	// InternalError indicates an internal JSON-RPC 2.0 error.
	ErrorCodeInternal ErrorCode = -32603
)

type RpcError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RpcMessage struct {
	Version string    `json:"version"`
	Id      uint64    `json:"id"`
	Method  string    `json:"method"`
	Params  []any     `json:"params"`
	Result  any       `json:"result"`
	Error   *RpcError `json:"error"`
}

func (r RpcMessage) IsCall() bool {
	return r.Id != 0 && r.Method != ""
}

func (r RpcMessage) IsNotify() bool {
	return r.Id == 0 && r.Method != ""
}

func (r RpcMessage) IsError() bool {
	return r.Error != nil
}

func (r RpcMessage) IsResult() bool {
	return r.Id != 0 && r.Result != nil
}

func MakeCall(id uint64, method string, params []any) *RpcMessage {
	return &RpcMessage{
		Version: "2.0",
		Id:      id,
		Method:  method,
		Params:  params,
	}
}

func MakeError(code ErrorCode, message string, data interface{}) *RpcMessage {
	return &RpcMessage{
		Version: "2.0",
		Error: &RpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

func MakeNotify(method string, params []any) *RpcMessage {
	return &RpcMessage{
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}

func MakeResult(id uint64, result any) *RpcMessage {
	return &RpcMessage{
		Version: "2.0",
		Id:      id,
		Result:  result,
	}
}
