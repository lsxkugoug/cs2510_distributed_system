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
 */
func main() {
	serverObj0 := server.GetServerObj(0)
	f, err := os.OpenFile("./testEntry/logs/server0.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	serverObj0.ProceedTest()
}
