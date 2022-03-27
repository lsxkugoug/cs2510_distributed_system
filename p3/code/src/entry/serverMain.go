package main

import (
	"code/src/config"
	"code/src/server"
	"fmt"
)

func main() {
	var serverId int
	fmt.Println("please input serverId, which is corresponding to the config file, indexed from 0 to ", len(config.Servers)-1)
	fmt.Scanf("%d", &serverId)
	serverObj := server.GetServerObj(serverId)
	if config.IndexOf(serverObj.ServerIpPort, config.Servers) == -1 {
		fmt.Println("config wrong!!!")
		return
	}
	serverObj.Proceed()
}
