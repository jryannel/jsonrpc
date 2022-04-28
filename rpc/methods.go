package rpc

import (
	"fmt"
	"sync"
)

type MethodHandle func(params []any) (any, error)

type Methods struct {
	mutex   sync.Mutex
	methods map[string]MethodHandle
}

func NewRegistry() *Methods {
	return &Methods{
		methods: make(map[string]MethodHandle),
	}
}

func (r *Methods) RegisterMethod(name string, handle MethodHandle) {
	r.mutex.Lock()
	r.methods[name] = handle
	r.mutex.Unlock()
}

func (r *Methods) UnregisterMethod(name string) {
	r.mutex.Lock()
	delete(r.methods, name)
	r.mutex.Unlock()
}

func (r *Methods) GetMethod(name string) MethodHandle {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.methods[name]
}

func (r *Methods) CallMethod(name string, params []any) (any, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.methods[name] == nil {
		return nil, fmt.Errorf("method %s not found", name)
	}
	return r.methods[name](params)
}
