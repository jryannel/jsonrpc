package main

import (
	"github.com/apigear-io/jsonrpc"
)

type Calc struct {
	Total int
}

func (c *Calc) Add(params []any) (any, error) {
	c.Total += int(params[0].(float64))
	return c.Total, nil
}

func (c *Calc) Clear(params []any) (any, error) {
	c.Total = 0
	return c.Total, nil
}

func main() {
	hub := jsonrpc.NewHub()
	calc := &Calc{}
	hub.RegisterMethod("calc.add", calc.Add)
	hub.RegisterMethod("calc.clear", calc.Clear)
	server := jsonrpc.NewHTTPServer()
	server.Router().Get("/ws", hub.HandleRequest)
	server.Start(":8080")
}
