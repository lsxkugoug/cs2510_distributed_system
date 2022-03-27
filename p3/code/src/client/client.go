package client

import (
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"net/rpc"
)

type Client struct {
	ServerIpPort []string //serverIp:serverPort
}

func GetClientObj() *Client {
	clientObj := &Client{ServerIpPort: config.Servers}
	return clientObj
}

func (client *Client) ClientSendMessage(args *rpc2.OperationArgs, reply *rpc2.OperationReply, serverIpPort string) *rpc2.OperationReply {
	args.IsLeader = false
	//fmt.Printf("the operation %v\n", args)
	conn, err := rpc.Dial("tcp", serverIpPort)
	if err != nil {
		//fmt.Println("rpc.Dial failed")
		return reply
	}
	err = conn.Call("Server.OperationHandler", args, reply)
	if err != nil {
		fmt.Printf("calling server %s failed", err)
	}
	//fmt.Printf("the value:%v, the vecotor:%v \n", reply.DataValue, reply.DataVector)
	return reply
}
