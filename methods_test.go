package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test register method
func TestRegisterMethod(t *testing.T) {
	registry := NewRegistry()
	isCalled := false
	handle := func(params []any) (any, error) {
		isCalled = true
		return nil, nil
	}
	registry.RegisterMethod("test", handle)
	handle = registry.GetMethod("test")
	assert.NotNil(t, handle)
	handle(nil)
	assert.True(t, isCalled)
}

// test unregister method
func TestUnregisterMethod(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterMethod("test", func(params []any) (any, error) {
		return nil, nil
	})
	registry.UnregisterMethod("test")
	assert.Nil(t, registry.GetMethod("test"))
}

// test call method
func TestCallMethod(t *testing.T) {
	registry := NewRegistry()
	isCalled := false
	registry.RegisterMethod("test", func(params []any) (any, error) {
		isCalled = true
		return "test", nil
	})
	result, err := registry.CallMethod("test", []any{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, result, "test")
	assert.True(t, isCalled)
}
