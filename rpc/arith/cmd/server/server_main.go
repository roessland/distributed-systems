package main

import "fmt"
import "net/rpc"
import "net"
import "log"
import "net/http"
import "github.com/roessland/distributed-systems/rpc/arith/server"

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Now listening at 127.0.0.1:1234")
	log.Fatal(http.Serve(listener, nil))
}
