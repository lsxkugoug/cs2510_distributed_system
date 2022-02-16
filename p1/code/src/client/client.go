package client

import (
	rpc2 "code/src/rpc"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

type Client struct {
	UserName     string
	UserGroup    string
	ServerIpPort string //serverIp:serverPort
	ClientPort   string // port
	ClientIpPort string // the IpPort app listen
}

func (client *Client) ServerSendMessageHandler(args *rpc2.ServerSendMessageArgs, reply *rpc2.ServerSendMessageReply) error {
	for _, msg := range args.Messages {
		fmt.Println(msg)
		log.Println(msg)
	}
	reply.Status = true
	return nil
}

func (client *Client) ClientSendMessage(rpcCall *rpc.Client, message string) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	message = "(" + timeStr + " " + ",group:" + client.UserGroup + ",name:" + client.UserName + "):" + message
	log.Println(message)

	args := rpc2.ClientSendMessageArgs{UserName: client.UserName, Message: message}
	reply := rpc2.ClientSendMessageReply{}
	err := rpcCall.Call("Server.ClientSendMessageHandler", &args, &reply)
	if err != nil {
		fmt.Printf("message %s failed", message)
	}
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

func (client *Client) getUnreadMsgs(conn *rpc.Client) {

	args := rpc2.ClientGetMessageArgs{UserName: client.UserName, UserGroup: client.UserGroup, ClientIpPort: client.ClientIpPort}
	reply := rpc2.ClientGetMessageReply{}
	err := conn.Call("Server.ClientGetMessageHandler", &args, &reply)
	if err != nil {
		fmt.Println("get unread messages wrong")
		return
	}

	for _, msg := range reply.UnreadMessages {
		fmt.Println(msg)
		log.Println(msg)
	}
}

// SetAndGetRpcCall:
// 1. register RPC
// 2. send ClientGetMessage RPC to get the unreaded messages, it can't be asynchronized
// 3. return RPCCALL
func (client *Client) SetAndGetConn() (*rpc.Client, error, net.Listener) {
	// 1. register RPC
	err := rpc.RegisterName("Client", client)
	if err != nil {
		fmt.Println("Client  register wrong")
		return nil, nil, nil
	}
	listener, err := net.Listen("tcp", ":"+client.ClientPort) // bind a port to receive calling
	if err != nil {
		fmt.Printf("listen %d failed", client.ClientPort)
		return nil, nil, nil
	}
	go rpc.Accept(listener) // build another thread to process RPC msg

	// 2. send ClientGetMessage RPC to get the unread messages, it can't be asynchronized
	conn, err := rpc.Dial("tcp", client.ServerIpPort)
	client.getUnreadMsgs(conn)
	if err != nil {
		fmt.Println("rpc.Dial failed")
		return nil, nil, nil
	}

	return conn, nil, listener

}
