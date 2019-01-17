package main

import (
	"learn/crawler/rpc/client"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {

	rpc.Register(client.DemoService{})
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("connect error %v\n", err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}

}
