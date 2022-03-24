package main

import "code/src/server"

/**
* @program: code
* @description:
* @author: Shixiang Long
* @create: 2022-03-09 22:41
**/
func main() {
	serverObj0 := server.GetServerObj(2)
	serverObj0.ProceedTest()
}
