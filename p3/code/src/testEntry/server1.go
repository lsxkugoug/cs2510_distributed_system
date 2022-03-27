package main

import (
	"code/src/server"
	"log"
	"os"
)

/**
* @program: code
* @description:
* @author: Shixiang Long
* @create: 2022-03-09 22:41
**/
func main() {
	f, err := os.OpenFile("./testEntry/logs/server1.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	serverObj0 := server.GetServerObj(1)
	serverObj0.ProceedTest()
}
