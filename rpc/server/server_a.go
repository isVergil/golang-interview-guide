package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
}

type ServiceA struct {
}

func (s *ServiceA) Add(args *Args, res *int) error {
	*res = args.X + args.Y
	return nil
}

func main() {
	service := new(ServiceA)
	rpc.Register(service) //注册 rpc
	//rpc.HandleHTTP()      //基于 http
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal("listen err:", err)
	}

	log.Println("RPC server started on :9091")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("conn err:", err)
		}
		//rpc.ServeConn(conn)
		//使用 json 协议
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
