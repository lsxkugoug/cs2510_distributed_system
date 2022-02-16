package main

import (
	clientPackage "code/src/client"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("./testEntry/test1/log/Bob.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	Bob := clientPackage.Client{UserName: "Bob", UserGroup: "1",
		ServerIpPort: clientPackage.GetOutboundIP("1234"), ClientPort: "1221",
		ClientIpPort: clientPackage.GetOutboundIP("1221")}
	conn, err, listener := Bob.SetAndGetConn()
	defer listener.Close()
	defer conn.Close()

	if err != nil {
		fmt.Println("Bob.SetAndGetRpcCall() wrong")
		return
	}
	// wait test process quit
	time.Sleep(4 * time.Second)

}
