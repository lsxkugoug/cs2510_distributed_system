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

//Design a test to show that your system only commits changes when quorum is achieved.
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
	args.IsLeader = false
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
	cmd4 := exec.Command("/bin/bash", "-c", "go run ./testEntry/server4.go")
	cmd4.Start()
	fmt.Println("=====================I am client, try to write key of server 0 and key: b value: 1======================")
	time.Sleep(2 * time.Second)
	client.ClientSendMessage(args, reply, config.Servers[0])

	//wait all of the 0 ~ 3 processed the key value

	time.Sleep(2 * time.Second)

	args.Key = "b"
	args.IsLeader = false
	args.Operation = 1
	fmt.Println("=====================I am client, try to read key of server 4, key: b============================")
	client.ClientSendMessage(args, reply, config.Servers[3])

	time.Sleep(10 * time.Second)

	fmt.Println("=================== logs server0 ~ server 4 ======================")
	fmt.Println("server0: ==================================")
	openTxt("./testEntry/logs/server0.txt")
	fmt.Println("server1: ==================================")
	openTxt("./testEntry/logs/server1.txt")
	fmt.Println("server2: ==================================")
	openTxt("./testEntry/logs/server2.txt")
	fmt.Println("server3: ==================================")
	openTxt("./testEntry/logs/server3.txt")
	fmt.Println("server4: ==================================")
	openTxt("./testEntry/logs/server4.txt")
	fmt.Printf("\n \n \n")
}
