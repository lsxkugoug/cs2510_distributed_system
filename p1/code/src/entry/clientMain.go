package main

import (
	"bufio"
	clientPKG "code/src/client"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	// 1. receive the data
	var serverIpPort string
	var clientPort string
	var name string
	var group string
	fmt.Println("please input your name:")
	fmt.Scanf("%s", &name)
	fmt.Println("please input your group:")
	fmt.Scanf("%s", &group)
	fmt.Println("please input <serverIp:severPort:>")
	fmt.Scanf("%s", &serverIpPort)
	fmt.Println("please input the port you want to use this app," +
		" if you want to open multiple apps on the same machine, the ports should be different" +
		"eg: 1200 ~ 1230")
	fmt.Scanf("%s", &clientPort)
	client := clientPKG.Client{UserName: name, UserGroup: group, ServerIpPort: serverIpPort,
		ClientIpPort: clientPKG.GetOutboundIP(clientPort), ClientPort: clientPort}
	log.SetOutput(ioutil.Discard)
	rpcCall, err, listener := client.SetAndGetConn()
	defer rpcCall.Close()
	defer listener.Close()
	if err != nil {
		fmt.Printf("get rpcCall wrong\n")
	}
	// process input message
	for {
		inputReader := bufio.NewReader(os.Stdin)
		msg, err := inputReader.ReadString('\n')
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		msg = "(" + timeStr + " " + ",group:" + client.UserGroup + ",name:" + client.UserName + "):" + msg
		if err != nil {
			fmt.Printf("input wrong\n")
		}
		go client.ClientSendMessage(rpcCall, msg)
		time.Sleep(10 * time.Microsecond)
	}
}
