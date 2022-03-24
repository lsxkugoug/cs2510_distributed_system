package main

import (
	clientPKG "code/src/client"
	"code/src/config"
	rpc2 "code/src/rpc"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

// normal
//① open 0 ~ 3 servers, but not 4
//② client new b = 1 on server 0
//③ server 4 online
//③ client update  b = 2 on server 1
//④ client retrieve all of the value of b in each servers, print the corresponding key, value and vector
func Test1(t *testing.T) {
	f, err := os.OpenFile("./test1.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	client := clientPKG.GetClientObj()
	defer f.Close()
	log.SetOutput(f)
	var args = &rpc2.OperationArgs{}
	var reply = &rpc2.OperationReply{}
	args.Key = "b"
	args.Value = "1"
	args.IsClient = true
	args.Operation = 0
	log.Printf("client new b = 1 to server0, in this time, 0 ~ 3 are online, but 4 offline\n")
	cmd0 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server0.go")
	cmd0.Start()
	cmd1 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server1.go")
	cmd1.Start()
	cmd2 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server2.go")
	cmd2.Start()
	cmd3 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server3.go")
	cmd3.Start()

	time.Sleep(2 * time.Second)
	client.ClientSendMessage(args, reply, config.Servers[0])

	//wait all of the 0 ~ 3 processed the key value
	time.Sleep(2 * time.Second)
	log.Printf("4 online now\n")
	cmd4 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server4.go")
	cmd4.Start()
	time.Sleep(2 * time.Second)

	log.Printf("client update b = 2 to server0, in this time, 0 ~ 4 are online\n")
	args.Key = "b"
	args.Value = "2"
	args.IsClient = true
	args.Operation = 0
	client.ClientSendMessage(args, reply, config.Servers[1])
	time.Sleep(2 * time.Second)
	log.Printf("client retrieve all of the servers b = ?\n")
	for i := 0; i <= 4; i++ {
		args.Key = "b"
		args.IsClient = true
		args.Operation = 1
		client.ClientSendMessage(args, reply, config.Servers[i])
		log.Printf("server%d:", i)
		log.Printf("the value:%v, the vecotor:%v \n", reply.DataValue, reply.DataVector)
		if len(reply.DataValue) > 1 {
			log.Println("conflict happens !!!!!")
		}
	}
	//cmd0.Process.Signal(syscall.SIGTERM)
	//cmd1.Process.Signal(syscall.SIGTERM)
	//cmd2.Process.Signal(syscall.SIGTERM)
	//cmd3.Process.Signal(syscall.SIGTERM)
	//cmd4.Process.Signal(syscall.SIGTERM)
	time.Sleep(10 * time.Second)

	fmt.Println("===================For test1======================")
	openTxt("./test1.txt")
	fmt.Printf("\n \n \n")
}

// conflict
//①  open 0 ~ 4 (all) servers
//② client new b = 1 on server 0    and  client new b = 2 on server 3
func Test2(t *testing.T) {
	f, err := os.OpenFile("./test2.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	client := clientPKG.GetClientObj()

	defer f.Close()
	log.SetOutput(f)
	var args1 = &rpc2.OperationArgs{}
	var reply1 = &rpc2.OperationReply{}
	args1.Key = "b"
	args1.Value = "1"
	args1.IsClient = true
	args1.Operation = 0

	var args2 = &rpc2.OperationArgs{}
	var reply2 = &rpc2.OperationReply{}
	args2.Key = "b"
	args2.Value = "2"
	args2.IsClient = true
	args2.Operation = 0
	log.Printf("client new b = 1 to server0, in this time, 0 ~ 3 are online, but 4 offline\n")
	cmd0 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server0.go")
	cmd0.Start()
	cmd1 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server1.go")
	cmd1.Start()
	cmd2 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server2.go")
	cmd2.Start()
	cmd3 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server3.go")
	cmd3.Start()
	cmd4 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server4.go")
	cmd4.Start()
	time.Sleep(2 * time.Second)
	log.Printf("client update b = 1 to server0, and send b = 2 to server 4 simutaneously, 0 ~ 4 are online\n")
	go client.ClientSendMessage(args1, reply1, config.Servers[0])
	go client.ClientSendMessage(args2, reply2, config.Servers[3])
	time.Sleep(2 * time.Second)
	log.Printf("client retrieve all of the servers b = ?\n")

	var args = &rpc2.OperationArgs{}
	var reply = &rpc2.OperationReply{}
	for i := 0; i <= 4; i++ {
		args.Key = "b"
		args.IsClient = true
		args.Operation = 1
		client.ClientSendMessage(args, reply, config.Servers[i])
		log.Printf("server%d:", i)
		log.Printf("the value:%v, the vecotor:%v \n", reply.DataValue, reply.DataVector)
		if len(reply.DataValue) > 1 {
			log.Println("conflict happens !!!!!")
		}
	}

	//cmd0.Process.Signal(syscall.SIGTERM)
	//cmd1.Process.Signal(syscall.SIGTERM)
	//cmd2.Process.Signal(syscall.SIGTERM)
	//cmd3.Process.Signal(syscall.SIGTERM)
	//cmd4.Process.Signal(syscall.SIGTERM)
	time.Sleep(10 * time.Second)
	fmt.Println("===================For test2======================")
	openTxt("./test2.txt")
	fmt.Printf("\n \n \n")
}

func Test3(t *testing.T) {
	ori := []int{0, 0, 0, 1, 0}
	received := []int{0, 0, 0, 0, 0}

	fmt.Println(isConflict(received, ori))
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
