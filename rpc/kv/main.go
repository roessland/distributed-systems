package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

//
// Common between client and server
//

const (
	OK       = "OK"
	ErrNoKey = "ErrNoKey"
)

type Err string

type PutArgs struct {
	Key   string
	Value string
}

type PutReply struct {
	Err Err
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Err   Err
	Value string
}

//
// Server
//

type KV struct {
	lock sync.Mutex
	data map[string]string
}

func (kv *KV) Get(args *GetArgs, reply *GetReply) error {
	kv.lock.Lock()
	log.Println("KV Get args", args)
	defer kv.lock.Unlock()
	if val, ok := kv.data[args.Key]; ok {
		reply.Value = val
		reply.Err = OK
	} else {
		reply.Value = ""
		reply.Err = ErrNoKey
	}
	log.Println("KV Get reply", reply)
	return nil
}

func (kv *KV) Put(args *PutArgs, reply *PutReply) error {
	kv.lock.Lock()
	log.Println("KV Put args", args)
	defer kv.lock.Unlock()
	kv.data[args.Key] = args.Value
	reply.Err = OK
	return nil
}

func Dial() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1235")
	if err != nil {
		log.Fatal("Dial:", err)
	}
	return client
}

func Get(key string) string {
	client := Dial()
	getArgs := GetArgs{key}
	getReply := GetReply{}
	err := client.Call("KV.Get", getArgs, &getReply)
	if err != nil {
		log.Fatal("Get call:", err)
	}
	return getReply.Value
}

func Put(key, val string) {
	client := Dial()
	putArgs := PutArgs{key, val}
	putReply := PutReply{}
	err := client.Call("KV.Put", putArgs, &putReply)
	if err != nil {
		log.Fatal("Put call:", err)
	}
}

// Clients

func server() {
	kv := new(KV)
	kv.data = map[string]string{}
	rpc.Register(kv)
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", "127.0.0.1:1235")
	if err != nil {
		log.Fatal("listen:", err)
	}
	go http.Serve(ln, nil)
}

func putClient(data map[string]string, wg *sync.WaitGroup) {
	for key, val := range data {
		Put(key, val)
	}
	wg.Done()
}

func readClient(keys []string, wg *sync.WaitGroup) {
	for _, key := range keys {
		Get(key)
	}
	wg.Done()
}

func main() {
	// important detail:
	// wait for net.Listen before continuing with client threads
	// so that the connections will succeed.
	server()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go putClient(map[string]string{
		"A": "a",
		"B": "b",
		"C": "c",
	}, &wg)
	go putClient(map[string]string{
		"A": "aa",
		"B": "bb",
		"C": "cc",
	}, &wg)
	go readClient([]string{"A", "B", "B", "B", "C", "D", "D", "A"}, &wg)
	go putClient(map[string]string{
		"B": "aaa",
		"C": "bbb",
		"D": "ccc",
	}, &wg)
	wg.Wait()
	fmt.Println("Done waiting for waitgroup")

}
