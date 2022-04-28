package rpc

import (
	"log"
	"sync"
)

type Connections struct {
	mutex       sync.Mutex
	connections map[*Connection]bool
}

func NewConnections() *Connections {
	return &Connections{
		connections: make(map[*Connection]bool),
	}
}

func (r *Connections) AddConnection(conn *Connection) {
	log.Println("Registering connection")
	r.mutex.Lock()
	r.connections[conn] = true
	r.mutex.Unlock()
}

func (r *Connections) RemoveConnection(conn *Connection) {
	log.Println("Unregister connection")
	r.mutex.Lock()
	delete(r.connections, conn)
	r.mutex.Unlock()
}

func (r *Connections) BroadcastMessage(msg *RpcMessage) {
	log.Println("Broadcasting message")
	r.mutex.Lock()
	for conn := range r.connections {
		conn.SendMessage(msg)
	}
	r.mutex.Unlock()
}

func (r *Connections) RemoveAllConnections() {
	log.Println("Closing all connections")
	for conn := range r.connections {
		r.RemoveConnection(conn)
	}
}
