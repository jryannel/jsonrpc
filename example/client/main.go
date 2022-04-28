package main

import (
	"fmt"

	"github.com/apigear-io/jsonrpc"
)

func main() {
	ws, err := jsonrpc.NewWebSocket("ws://localhost:8080/ws")
	if err != nil {
		panic(err)
	}
	client := jsonrpc.NewRpcClient(ws)
	for i := 0; i < 100000; i++ {
		result, err := client.Call("calc.add", []any{i})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", result)
	}
	result, err := client.Call("calc.clear", []any{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", result)
}
