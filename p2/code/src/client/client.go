package client

import (
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"net"
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
	args.IsClient = true
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
	if len(reply.DataValue) > 1 {
		fmt.Println("conflict happens !!!!!")
	}
	return reply
}
func GetOutboundIP(port string) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String() + ":" + port
}
