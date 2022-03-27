package main

import (
	clientPKG "code/src/client"
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// 1. receive the data
	client := clientPKG.GetClientObj()
	var serverId int
	var operation int
	var key string
	var value string
	var args = &rpc2.OperationArgs{}
	var reply = &rpc2.OperationReply{}
	log.SetOutput(os.Stdout)
	// process input message
	for {
		fmt.Println("please input the operation you want, 0:update or new 1: retrieve data")
		fmt.Scanf("%d", &operation)
		args.Operation = operation
		fmt.Printf("please input the serverId you want to make operation, you can input 0 ~ %d\n", len(config.Servers)-1)
		fmt.Scanf("%d", &serverId)
		// if update or new
		if operation == 0 {
			fmt.Println("please input the key you want update or create")
			fmt.Scanf("%s", &key)
			args.Key = key
			fmt.Println("please input the corresponding value")
			fmt.Scanf("%s", &value)
			args.Value = value
			client.ClientSendMessage(args, reply, config.Servers[serverId])
		}
		if operation == 1 {
			fmt.Println("please input the key you want retrieve")
			fmt.Scanf("%s", &key)
			args.Key = key
			client.ClientSendMessage(args, reply, config.Servers[serverId])
			log.Printf("the value:%s,", reply.Value)
		}

		time.Sleep(10 * time.Microsecond)
	}
}
