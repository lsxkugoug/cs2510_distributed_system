package main

import (
	clientPackage "code/src/client"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// important, test_test.go is under src, therefore, the relative path is "./testEntry/test1/log.Alice"
	f, err := os.OpenFile("./testEntry/test1/log/Alice.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	Alice := clientPackage.Client{UserName: "Alice", UserGroup: "1",
		ServerIpPort: clientPackage.GetOutboundIP("1234"), ClientPort: "1220",
		ClientIpPort: clientPackage.GetOutboundIP("1220")}
	conn, err, listener := Alice.SetAndGetConn()
	defer conn.Close()
	defer listener.Close()
	if err != nil {
		fmt.Println("Alice.SetAndGetRpcCall() wrong")
		return
	}
	Alice.ClientSendMessage(conn, "I am Alice, nice to meet you")
	// wait test process quit
	time.Sleep(8 * time.Second)
}
