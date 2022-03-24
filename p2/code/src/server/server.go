package server

import (
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Server struct {
	mu           sync.Mutex
	ServerId     int
	KeyValue     map[string][]string // key value
	KeyVector    map[string][][]int  // key vector
	LocalClock   []int
	ServerIpPort string
}

func GetServerObj(serverId int) Server {
	serverObj := Server{
		ServerId:     serverId,
		ServerIpPort: config.Servers[serverId],
		KeyValue:     map[string][]string{},
		KeyVector:    map[string][][]int{},
		LocalClock:   make([]int, len(config.Servers)),
	}
	fmt.Printf("the parameter of server %v\n", serverObj)
	return serverObj
}
func isConflict(ori []int, received []int) bool {
	flg := 0
	// ori > rec, flg = 1, ori < rec, flg = -1
	// [1,0,1], [0,1,0] is a conflict situation
	for idx, v := range ori {
		if received[idx] > v {
			if flg < 0 {
				return true
			}
			flg = 1
		}
		if received[idx] == v {
			continue
		}
		if received[idx] < v {
			if flg > 0 {
				return true
			}
			flg = -1
		}
	}
	return false
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

func (server *Server) OperationHandler(args *rpc2.OperationArgs, reply *rpc2.OperationReply) error {
	server.mu.Lock()
	defer server.mu.Unlock()
	//fmt.Printf("get the args %v, the current data%v\n", args, server.KeyVector[args.Key])
	if args.IsClient {
		args.Vector = make([]int, len(config.Servers))
		// receive a new update or new. add my local clock, if the operation is update or new
		if args.Operation == 0 {
			server.LocalClock[server.ServerId] += 1
		}
	}

	// justify whether conflict, if client send, never conflict
	confligFlag := isConflict(server.LocalClock, args.Vector)

	// update local vector
	//fmt.Println(server.LocalClock, args.Vector)
	if !confligFlag {
		for idx, v := range server.LocalClock {
			if args.Vector[idx] > v {
				server.LocalClock[idx] = args.Vector[idx]
			}
		}
	}

	if args.Operation == 0 { // update or new data
		_, ok := server.KeyValue[args.Key]
		if !ok { // if ok not in the keyValue map
			server.KeyValue[args.Key] = append(server.KeyValue[args.Key], args.Value)
			server.KeyVector[args.Key] = append(server.KeyVector[args.Key], server.LocalClock)
			//fmt.Printf("new data %v, %v\n", args, server.KeyVector[args.Key])
		} else {
			if len(server.KeyValue[args.Key]) > 1 { // means already conflict
				server.KeyValue[args.Key] = append(server.KeyValue[args.Key], args.Value)
				server.KeyVector[args.Key] = append(server.KeyVector[args.Key], server.LocalClock)
				//fmt.Printf("update data and len > 1, conflict %v\n", args)
			} else {
				// check whether the vector conflict
				if confligFlag {
					server.KeyValue[args.Key] = append(server.KeyValue[args.Key], args.Value)
					server.KeyVector[args.Key] = append(server.KeyVector[args.Key], args.Vector)
				} else {
					server.KeyValue[args.Key][0] = args.Value
					server.KeyVector[args.Key][0] = server.LocalClock
					//fmt.Printf("update data rewrite %v\n", args)
				}
			}
		}
	} else { // retrieve the data
		_, ok := server.KeyValue[args.Key]
		if !ok {
			return nil
		}
		reply.Operation = args.Operation
		reply.DataValue = server.KeyValue[args.Key]
		reply.DataVector = server.KeyVector[args.Key]
	}

	// send the information to other servers
	if !args.IsClient {
		return nil
	}

	// try to broadcast new update or new. add my local clock, if the operation is update or new
	if args.Operation == 0 {
		server.LocalClock[server.ServerId] += 1
	}
	// only operation is 0, we need to broadcast it
	if args.Operation == 0 {
		for idx, _ := range config.Servers {
			args.IsClient = false
			args.Vector = server.LocalClock
			if idx == server.ServerId {
				continue
			}
			go sendMsg(args, idx)
		}
	}
	return nil
}

func GetOutboundIP(port string) string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String() + ":" + port
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
