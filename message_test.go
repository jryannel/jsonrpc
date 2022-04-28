package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCall(t *testing.T) {
	msg := MakeCall(1, "test", []any{1, 2, 3})
	assert.Equal(t, msg.Version, "2.0")
	assert.Equal(t, msg.Id, uint64(1))
	assert.Equal(t, msg.Method, "test")
	assert.Equal(t, msg.Params, []any{1, 2, 3})
	assert.Nil(t, msg.Error, nil)
	assert.Nil(t, msg.Result, nil)
	assert.True(t, msg.IsCall())
	assert.False(t, msg.IsNotify())
	assert.False(t, msg.IsError())
	assert.False(t, msg.IsResult())
}

// test make notify
func TestMakeNotify(t *testing.T) {
	msg := MakeNotify("test", []any{1, 2, 3})
	assert.Equal(t, msg.Version, "2.0")
	assert.Equal(t, msg.Method, "test")
	assert.Equal(t, msg.Params, []any{1, 2, 3})
	assert.Equal(t, msg.Id, uint64(0))
	assert.Nil(t, msg.Error, nil)
	assert.Nil(t, msg.Result, nil)
	assert.False(t, msg.IsCall())
	assert.True(t, msg.IsNotify())
	assert.False(t, msg.IsError())
	assert.False(t, msg.IsResult())
}

// test make error
func TestMakeError(t *testing.T) {
	msg := MakeError(ErrorCodeParse, "test", "data")
	assert.Equal(t, msg.Version, "2.0")
	assert.Equal(t, msg.Error.Code, ErrorCodeParse)
	assert.Equal(t, msg.Error.Message, "test")
	assert.Equal(t, msg.Error.Data, "data")
	assert.Nil(t, msg.Result, nil)
	assert.False(t, msg.IsCall())
	assert.False(t, msg.IsNotify())
	assert.True(t, msg.IsError())
	assert.False(t, msg.IsResult())
}

// test make result
func TestMakeResult(t *testing.T) {
	msg := MakeResult(1, "test")
	assert.Equal(t, msg.Version, "2.0")
	assert.Equal(t, msg.Id, uint64(1))
	assert.Equal(t, msg.Result, "test")
	assert.Nil(t, msg.Error, nil)
	assert.Nil(t, msg.Params, nil)
	assert.False(t, msg.IsCall())
	assert.False(t, msg.IsNotify())
	assert.False(t, msg.IsError())
	assert.True(t, msg.IsResult())
}
