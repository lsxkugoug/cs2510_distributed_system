package server

import (
	rpc2 "code/src/rpc"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Server struct {
	mu           sync.Mutex
	NameGroup    map[string]string   // k: name, v: group
	NameIps      map[string]string   // k: name, v: message
	GroupMessage map[string][]string // k: group	v: list of message
	NameRead     map[string]int      // k: name  v: readIdx
}

func (server *Server) ClientGetMessageHandler(args *rpc2.ClientGetMessageArgs, reply *rpc2.ClientGetMessageReply) error {
	server.mu.Lock()
	defer server.mu.Unlock()
	_, ok := server.NameGroup[args.UserName]
	// 1. if it is the first time this user try to send the message, register it
	if !ok {
		fmt.Println(args.UserName)
		server.NameGroup[args.UserName] = args.UserGroup
		server.NameIps[args.UserName] = args.ClientIpPort
		server.NameRead[args.UserName] = -1
	}
	readIdx := server.NameRead[args.UserName]
	replyMsgs := server.GroupMessage[args.UserGroup][readIdx+1:]
	server.NameRead[args.UserName] = len(server.GroupMessage[args.UserGroup]) - 1
	reply.UnreadMessages = replyMsgs
	return nil
}

func (server *Server) ClientSendMessageHandler(args *rpc2.ClientSendMessageArgs, reply *rpc2.ClientSendMessageReply) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	// 1. add the message to the group map
	userGroup := server.NameGroup[args.UserName]
	server.GroupMessage[userGroup] = append(server.GroupMessage[userGroup], args.Message)
	server.NameRead[args.UserName] = len(server.GroupMessage[userGroup]) - 1

	// flooding the message
	for name, group := range server.NameGroup {
		if group == userGroup {
			if name == args.UserName {
				continue
			}
			readIdx := server.NameRead[name]
			unreadedMsg := server.GroupMessage[group][readIdx+1:]
			go server.sendMsgToClient(name, server.NameIps[name], unreadedMsg, len(server.GroupMessage[group])-1)
		}
	}
	return nil
}

func (server *Server) sendMsgToClient(userName string, clientIpPort string, messages []string, readIndx int) {
	rpcCall, err := rpc.Dial("tcp", clientIpPort)
	if err != nil {
		fmt.Printf("server sendMsgToClient rpc.Dial clientIpPort%v\n", clientIpPort)
		return
	}
	args := rpc2.ServerSendMessageArgs{Messages: messages}
	reply := rpc2.ServerSendMessageReply{}
	err = rpcCall.Call("Client.ServerSendMessageHandler", &args, &reply)
	if err != nil || !reply.Status {
		fmt.Printf("server send messages %v wrong\n", messages)
		return
	}
	// if successfull send the client, let client's server.NameRead = readIndx
	server.mu.Lock()
	if server.NameRead[userName] <= readIndx {
		server.NameRead[userName] = readIndx
	}
	server.mu.Unlock()

}

func GetOutboundIP(port string) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String() + ":" + port
}

// Proceed :test and normal entry
func (server *Server) Proceed() {
	port := "1234"

	// register RPC
	err := rpc.RegisterName("Server", server)
	if err != nil {
		fmt.Println("Server register failed")
	}
	listener, err := net.Listen("tcp", ":"+port) // bind a port to receive calling
	if err != nil {
		fmt.Println("Listen port failed")
		return
	}
	fmt.Println(GetOutboundIP(port))
	defer listener.Close()
	rpc.Accept(listener) // loop forever
}
