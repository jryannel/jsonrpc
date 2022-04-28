package rpc

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test protocol notify
func TestProtocolNotify(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	p.Notify("test", []any{1, 2, 3})
	assert.Equal(t, len(s.Messages), 1)
	assert.Equal(t, s.Messages[0].Method, "test")
}

func TestFailingHandleNotify(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeNotify("test", []any{1, 2, 3})
	p.HandleMessage(msg)
	assert.Equal(t, 1, len(s.Messages))
	assert.Equal(t, ErrorCodeMethodNotFound, s.Messages[0].Error.Code)
}

// test protocol call
func TestProtocolCall(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	wg := sync.WaitGroup{}
	wg.Add(1)
	id := p.NextSeq()
	go func() {
		// send a result message
		// TODO: difficult to test as we don't know the id
		msg := MakeResult(id, "test")
		p.HandleMessage(msg)
		wg.Done()
	}()
	// Call will block. We need to wait for the call to be handled.
	result, err := p.CallWithId(id, "test", []any{1, 2, 3})
	assert.Equal(t, len(s.Messages), 1)
	assert.Equal(t, s.Messages[0].Id, uint64(1))
	assert.Equal(t, s.Messages[0].Method, "test")
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, result, "test")
}

// test protocol error
func TestProtocolError(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	p.Error(ErrorCodeParse, "test")
	assert.Equal(t, len(s.Messages), 1)
	assert.Equal(t, s.Messages[0].Error.Code, ErrorCodeParse)
	assert.Equal(t, s.Messages[0].Error.Message, "test")
}

func TestHandleNotify(t *testing.T) {
	r := NewRegistry()
	isCalled := false
	r.RegisterMethod("test", func(params []any) (any, error) {
		isCalled = true
		return nil, nil
	})
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeNotify("test", []any{1, 2, 3})
	p.HandleMessage(msg)
	assert.True(t, isCalled)
}

func TestHandleError(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeError(ErrorCodeParse, "test", nil)
	p.HandleMessage(msg)
	assert.Equal(t, len(s.Messages), 0)
}

func TestHandleCallError(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeError(ErrorCodeInternal, "test", []any{1, 2, 3})
	msg.Id = 1
	p.HandleMessage(msg)
	assert.Equal(t, len(s.Messages), 0)
}

// test handle call
func TestHandleCall(t *testing.T) {
	r := NewRegistry()
	isCalled := false
	r.RegisterMethod("test", func(params []any) (any, error) {
		isCalled = true
		return "test", nil
	})
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeCall(1, "test", []any{1, 2, 3})
	p.HandleMessage(msg)
	assert.True(t, isCalled)
	assert.Equal(t, len(s.Messages), 1)
	assert.Equal(t, s.Messages[0].Id, uint64(1))
	assert.Equal(t, s.Messages[0].Result, "test")
}

// test handle result
func TestHandleResult(t *testing.T) {
	r := NewRegistry()
	s := NewMockMessageSender()
	p := NewProtocol(s, r)
	msg := MakeResult(1, "test")
	p.HandleMessage(msg)
	assert.Equal(t, len(s.Messages), 0)
}
