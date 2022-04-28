package jsonrpc

import (
	"sync"
)

// mock message sender
type MockMessageSender struct {
	Messages []*RpcMessage
	mutex    sync.Mutex
}

func NewMockMessageSender() *MockMessageSender {
	return &MockMessageSender{
		Messages: make([]*RpcMessage, 0),
	}
}

func (m *MockMessageSender) SendMessage(msg *RpcMessage) error {
	m.mutex.Lock()
	m.Messages = append(m.Messages, msg)
	m.mutex.Unlock()
	return nil
}
