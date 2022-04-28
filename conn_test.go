package jsonrpc

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HttpToWsAddr(strAddr string) string {
	urlAddr, err := url.Parse(strAddr)
	if err != nil {
		return strAddr
	}
	urlAddr.Scheme = "ws"
	urlAddr.Path = "/ws"
	return urlAddr.String()
}

func makeTestClient(addr string) (*RpcClient, error) {
	ws, err := NewWebSocket(addr)
	if err != nil {
		return nil, err
	}
	return NewRpcClient(ws), nil
}

func makeTestHub() (*Hub, *httptest.Server) {
	hub := NewHub()
	handler := NewRouter()
	handler.Get("/ws", hub.HandleRequest)
	ts := httptest.NewServer(handler)
	return hub, ts
}

func TestNewConnection(t *testing.T) {
	// make server
	hub, ts := makeTestHub()
	// make client
	client, err := makeTestClient(HttpToWsAddr(ts.URL))
	assert.NoError(t, err)

	defer func() {
		hub.RemoveAllConnections()
		client.Close()
		ts.Close()
	}()

	hub.RegisterMethod("test", func(args []any) (any, error) {
		return "test", nil
	})
	result, err := client.Call("test", []any{})
	assert.NoError(t, err)
	assert.Equal(t, "test", result)

	done := make(chan bool)
	client.RegisterMethod("test2", func(args []any) (any, error) {
		done <- true
		return nil, nil
	})
	hub.Notify("test2", []any{})
	isCalled := <-done
	assert.True(t, isCalled)
}
