package main

import (
	client2 "learn/test/rpc/client"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	client := jsonrpc.NewClient(conn)

	var result float64
	err = client.Call("DemoService.Div", client2.Args{1, 3}, &result)

	if err != nil {
		log.Printf("call service error %v", err)
	} else {
		log.Println("call success result is ", result)
	}

	err = client.Call("DemoService.Div", client2.Args{1, 0}, &result)

	if err != nil {

		log.Printf("call service error %v", err)
	}
}
