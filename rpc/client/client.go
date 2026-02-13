package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
}

func main() {
	// 建立 HTTP 连接
	// client, err := rpc.Dial("tcp", "127.0.0.1:9091")
	// 基于 json 协议
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 同步调用
	args := &Args{
		X: 10,
		Y: 20,
	}
	var ret1 int
	err = client.Call("ServiceA.Add", args, &ret1)
	if err != nil {
		log.Fatal("ServiceA.Add err:", err)
	}
	fmt.Printf("ServiceA.Add: %d + %d = %d \n", args.X, args.Y, ret1)

	// 异步调用
	for i := 0; i < 1000; i++ {
		go func() {
			var ret2 int
			conn1, _ := net.Dial("tcp", "127.0.0.1:9091")
			client1 := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn1))
			async1 := client1.Go("ServiceA.Add", args, &ret2, nil)
			res2Call := <-async1.Done
			fmt.Println(res2Call.Error)
			fmt.Println(ret2)
		}()
	}

}
