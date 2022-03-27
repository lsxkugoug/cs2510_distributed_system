package server

import (
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Server struct {
	mu           sync.Mutex
	ServerId     int
	KeyValue     map[string]string // key value
	UnCommit     map[string]string // key value
	state        int               //0 leader, 1 normal server
	leader       int               // 0  len(nodes) - 1, initial value = -1
	ServerIpPort string
	MyTerm       int // the order I make the election

}

func GetServerObj(serverId int) Server {
	serverObj := Server{
		ServerId:     serverId,
		ServerIpPort: config.Servers[serverId],
		KeyValue:     map[string]string{},
		UnCommit:     map[string]string{},
		state:        1,
		leader:       -1,
		MyTerm:       0,
	}
	fmt.Printf("the parameter of server %v\n", serverObj)
	return serverObj
}

func sendMsg(args *rpc2.OperationArgs, serverId int) {
	fmt.Printf("broatcast operation:%v, to serverId: %d\n", args, serverId)
	reply := &rpc2.OperationReply{}
	conn, err := rpc.Dial("tcp", config.Servers[serverId])
	if err != nil {
		//fmt.Println("rpc.Dial failed")
		return
	}
	err = conn.Call("Server.OperationHandler", args, reply)
	if err != nil {
		fmt.Printf("calling server %s failed", err)
	}
	return
}
func (server *Server) leaderElection() error {
	server.mu.Lock()
	defer server.mu.Unlock()
	log.Println("I am ", server.ServerId, " try to make leader election")

	// send the request to all the nodes
	for serverId, _ := range config.Servers {
		if serverId <= server.ServerId {
			continue
		}
		args := &rpc2.LeaderElectionArgs{Type: 0, SenderId: server.ServerId, SenderTerm: server.MyTerm}
		reply := &rpc2.LeaderElectionReply{}
		conn, err := rpc.Dial("tcp", config.Servers[serverId])
		if err != nil {
			log.Println("I am", server.ServerId, "and try to send leader election msg to", serverId, "but failed")
			continue
		}
		err = conn.Call("Server.LeaderElectionHandler", args, reply)
		if err != nil {
			log.Printf("calling server %s failed", err)
		}
		if reply.IsOk {
			return nil
		}
	}
	// if doesn't get isOk, flooding the "I am won"
	server.leader = server.ServerId
	log.Println("I am leader, I am ", server.leader)
	for serverId, _ := range config.Servers {
		if serverId == server.ServerId {
			continue
		}
		args := &rpc2.LeaderElectionArgs{Type: 1, SenderId: server.ServerId}
		reply := &rpc2.LeaderElectionReply{}
		conn, err := rpc.Dial("tcp", config.Servers[serverId])
		if err != nil {
			log.Println("I am", server.ServerId, "try to send \"I am won msg to\"", serverId, "but failed")
			continue
		}
		err = conn.Call("Server.LeaderElectionHandler", args, reply)
		if err != nil {
			fmt.Printf("calling server %s failed", err)
			continue
		}
	}
	return nil
}
func (server *Server) OperationHandler(args *rpc2.OperationArgs, reply *rpc2.OperationReply) error {
	// read operation
	if args.IsRead {
		if v, ok := server.KeyValue[args.Key]; ok {
			reply.Value = v
		} else {
			reply.Value = ""
		}
		log.Println("I am", server.ServerId, "get the read request from another and return", reply.Value)
		return nil
	}
	if args.Operation == 1 {
		for {
			var countMap = make(map[string]int) // key count
			for serverId, _ := range config.Servers {
				if serverId == server.ServerId {
					continue
				}
				argsRead := &rpc2.OperationArgs{Operation: 1, IsRead: true, Key: args.Key}
				replyRead := &rpc2.OperationReply{}
				conn, err := rpc.Dial("tcp", config.Servers[serverId])
				if err != nil {
					log.Println("I am", server.ServerId, "try to read data from", serverId, "but failed")
					continue
				}
				err = conn.Call("Server.OperationHandler", argsRead, replyRead)
				if err != nil {
					fmt.Printf("calling server %s failed", err)
					continue
				}
				if _, ok := countMap[replyRead.Value]; !ok {
					countMap[replyRead.Value] = 0
				} else {
					countMap[replyRead.Value] += 1
				}
			}
			for k, v := range countMap {
				if v >= config.Nr {
					log.Printf("I am %d,read data from other replicas and the number of %d > Nr(%d) and return it to the client", server.ServerId, v, config.Nr)
					reply.Value = k
					return nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
	// forwarding the request to the leader, if server.leader == -1, start leader election
	if server.leader == -1 {
		log.Println("I am", server.ServerId, "find my leader variable == -1, start leader election")
		server.mu.Lock()
		server.MyTerm++
		server.mu.Unlock()
		server.leaderElection()
	}
	// waiting for server.leader!=-1
	for {
		server.mu.Lock()
		if server.leader != -1 {
			server.mu.Unlock()
			break
		}
		server.mu.Unlock()
	}

	// if I am not leader and the message is come from client, forward it to leader
	if server.leader != server.ServerId && !args.IsLeader {
		// if the request is update_new, should forward to the leader
		conn, err := rpc.Dial("tcp", config.Servers[server.leader])
		if err != nil {
			log.Println("I am", server.ServerId, "try to call leader", server.leader, "but failed")
			server.mu.Lock()
			server.MyTerm++
			server.mu.Unlock()
			server.leaderElection()
		}
		// wait for leader change
		for {
			conn, err = rpc.Dial("tcp", config.Servers[server.leader])
			if err == nil {
				break
			}
		}
		err = conn.Call("Server.OperationHandler", args, reply)
		if err != nil {
			fmt.Printf("calling server %s failed", err)
			return nil
		}
		return nil
	}

	// the later part is update or read operation
	// if I am leader
	if server.ServerId == server.leader {
		isReadyCount := 1
		// get quorum
		log.Println("I am leader:", server.ServerId, "try to get quorum ")

		for serverId, _ := range config.Servers {
			if serverId == server.ServerId {
				continue
			}
			argsLeader := &rpc2.OperationArgs{Operation: 0,
				IsLeader: true,
				Key:      args.Key,
				Value:    args.Value,
				IsCommit: false,
			}
			replyLeader := &rpc2.OperationReply{}
			conn, err := rpc.Dial("tcp", config.Servers[serverId])
			if err != nil {
				log.Println("I am", server.ServerId, "try to get quorum data from", serverId, "but failed")
				continue
			}
			err = conn.Call("Server.OperationHandler", argsLeader, replyLeader)
			if err != nil {
				fmt.Printf("calling server %s failed", err)
				continue
			}
			if replyLeader.IsReady {
				isReadyCount++
			}
		}
		// if can't get quorum, stop the operation
		if isReadyCount < config.Nw {
			return nil
		}
		// send commit
		log.Println("I am leader:", server.ServerId, ", I have already get quorum and try to send commit")

		server.KeyValue[args.Key] = args.Value
		for serverId, _ := range config.Servers {
			if serverId == server.ServerId {
				continue
			}
			argsLeader := &rpc2.OperationArgs{Operation: 0,
				IsLeader: true,
				Key:      args.Key,
				Value:    args.Value,
				IsCommit: true,
			}
			replyLeader := &rpc2.OperationReply{}
			conn, err := rpc.Dial("tcp", config.Servers[serverId])
			if err != nil {
				log.Println("I am leader:", server.ServerId, "try to send commit", serverId, "but failed")
				continue
			}
			err = conn.Call("Server.OperationHandler", argsLeader, replyLeader)
			if err != nil {
				fmt.Printf("calling server %s failed", err)
				continue
			}
		}
	}
	// if I am not leader
	if server.ServerId != server.leader {
		if args.IsCommit == false {
			log.Println("I am:", server.ServerId, "stored the", args.Key, args.Value, "into uncommit datastructure")
			server.UnCommit[args.Key] = args.Value
			reply.IsReady = true
			return nil
		} else {
			log.Println("I am:", server.ServerId, "get the commit instruction and apply key =", args.Key, "and value", args.Value)
			if _, ok := server.UnCommit[args.Key]; ok {
				delete(server.UnCommit, args.Key)
				server.KeyValue[args.Key] = args.Value
			}
		}
	}
	return nil
}

func (server *Server) LeaderElectionHandler(args *rpc2.LeaderElectionArgs, reply *rpc2.LeaderElectionReply) error {
	server.mu.Lock()
	if args.Type == 1 { // if the msg is 'I am won'
		server.leader = args.SenderId
		log.Println("Leader is ", server.leader)
	} else {
		if args.SenderId > server.ServerId {
			reply.IsOk = false
		} else {
			reply.IsOk = true
		}
	}

	if args.SenderTerm > server.MyTerm {
		server.MyTerm = args.SenderTerm
		log.Println("I am ", server.ServerId, ", receive ", args.SenderId, "'s leader election message, so participate the leader election")
		server.mu.Unlock()
		go server.leaderElection()
		server.mu.Lock()
	}
	server.mu.Unlock()
	return nil
}

// Proceed :test and normal entry
func (server *Server) Proceed() {
	// register RPC
	err := rpc.RegisterName("Server", server)
	if err != nil {
		fmt.Println("Server register failed")
	}
	//listener, err := net.Listen("tcp", ":"+port) // bind a port to receive calling
	listener, err := net.Listen("tcp", server.ServerIpPort) // bind a port to receive calling
	if err != nil {
		fmt.Println("Listen port failed")
		return
	}
	defer listener.Close()
	rpc.Accept(listener) // loop forever
}

// Proceed :test and normal entry
func (server *Server) ProceedTest() {
	// register RPC
	err := rpc.RegisterName("Server", server)
	if err != nil {
		fmt.Println("Server register failed")
	}
	//listener, err := net.Listen("tcp", ":"+port) // bind a port to receive calling
	listener, err := net.Listen("tcp", server.ServerIpPort) // bind a port to receive calling

	if err != nil {
		fmt.Println("Listen port failed")
		return
	}
	defer listener.Close()
	go rpc.Accept(listener) // loop forever
	time.Sleep(15 * time.Second)
}
