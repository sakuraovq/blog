package rpcsupport

import (
	"gopkg.in/olivere/elastic.v3"
	"learn/crawler/engine"
	"learn/crawler/persist"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (saver ItemSaverService) Saver(item engine.Item, result *string) error {

	error := persist.Save(item, saver.Client, saver.Index)
	if error != nil {
		return error
	}
	log.Printf("item saved %v", item)
	*result = "ok"
	return nil
}

// 初始化 rpcClient
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, err
	}
	client := jsonrpc.NewClient(conn)
	return client, nil
}

func RpcServer(host string, service interface{}) error {

	rpc.Register(service)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		return err
	}
	log.Print("server run host ", host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("connect error %v\n", err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
