package rpc

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

type MethodCaller interface {
	CallMethod(name string, params []any) (any, error)
}

type MessageSender interface {
	SendMessage(msg *RpcMessage) error
}

type PendingCall struct {
	Message *RpcMessage
	Done    chan *RpcMessage
}

func NewPendingCall(msg *RpcMessage) *PendingCall {
	return &PendingCall{
		Message: msg,
		Done:    make(chan *RpcMessage),
	}
}

type Protocol struct {
	seq     uint64
	sender  MessageSender
	caller  MethodCaller
	mutex   sync.Mutex
	pending map[uint64]*PendingCall
}

func NewProtocol(sender MessageSender, caller MethodCaller) *Protocol {
	return &Protocol{
		sender:  sender,
		caller:  caller,
		pending: make(map[uint64]*PendingCall),
	}
}

func (p *Protocol) NextSeq() uint64 {
	atomic.AddUint64(&p.seq, 1)
	return p.seq
}

func (p *Protocol) takePending(id uint64) (*PendingCall, bool) {
	p.mutex.Lock()
	call, ok := p.pending[id]
	delete(p.pending, id)
	p.mutex.Unlock()
	if !ok {
		return nil, false
	}
	return call, true
}

func (p *Protocol) HandleMessage(msg *RpcMessage) {
	if msg.IsCall() {
		p.HandleCall(msg)
	} else if msg.IsNotify() {
		p.HandleNotify(msg)
	} else if msg.IsResult() {
		p.HandleResult(msg)
	} else if msg.IsError() {
		p.HandleError(msg)
	}
}

func (p *Protocol) HandleCall(msg *RpcMessage) {
	result, err := p.caller.CallMethod(msg.Method, msg.Params)
	if err != nil {
		p.SendMessage(MakeError(ErrorCodeMethodNotFound, err.Error(), err))
		return
	}
	p.SendMessage(MakeResult(msg.Id, result))
}

func (p *Protocol) HandleNotify(msg *RpcMessage) {
	_, err := p.caller.CallMethod(msg.Method, msg.Params)
	if err != nil {
		p.SendMessage(MakeError(ErrorCodeMethodNotFound, err.Error(), err))
	}
}

func (p *Protocol) HandleResult(msg *RpcMessage) {
	call, ok := p.takePending(msg.Id)
	if !ok {
		log.Printf("rpc: unknown result id: %d", msg.Id)
		return
	}
	call.Done <- msg
}

func (p *Protocol) HandleError(msg *RpcMessage) {
	if msg.Id != 0 {
		call, ok := p.takePending(msg.Id)
		if ok {
			call.Done <- msg
		}
	}
	log.Printf("rpc: error: %d: %s", msg.Error.Code, msg.Error.Message)
}

func (p *Protocol) SendMessage(msg *RpcMessage) error {
	return p.sender.SendMessage(msg)
}

func (p *Protocol) CallWithId(id uint64, method string, params []any) (any, error) {
	msg := MakeCall(id, method, params)
	call := NewPendingCall(msg)
	p.mutex.Lock()
	p.pending[msg.Id] = call
	p.mutex.Unlock()
	err := p.SendMessage(msg)
	if err != nil {
		return nil, err
	}
	// block until call is done
	result := <-call.Done
	close(call.Done)
	if result.Error != nil {
		return nil, fmt.Errorf("jsonrpc error: %d: %s", result.Error.Code, result.Error.Message)
	}
	return result.Result, nil
}

func (p *Protocol) Call(method string, params []any) (any, error) {
	return p.CallWithId(p.NextSeq(), method, params)
}

func (p *Protocol) Notify(method string, params []any) error {
	msg := MakeNotify(method, params)
	return p.SendMessage(msg)
}

func (p *Protocol) Error(code ErrorCode, message string) error {
	msg := MakeError(code, message, nil)
	return p.SendMessage(msg)
}
