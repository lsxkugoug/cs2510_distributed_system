package main

import (
	clientPackage "code/src/client"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	f, err := os.OpenFile("./testEntry/test2/log/Dough.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	Doug := clientPackage.Client{UserName: "Doug", UserGroup: "2",
		ServerIpPort: clientPackage.GetOutboundIP("1234"), ClientPort: "1223",
		ClientIpPort: clientPackage.GetOutboundIP("1223")}
	conn, err, listener := Doug.SetAndGetConn()
	defer listener.Close()
	defer conn.Close()
	if err != nil {
		fmt.Println("Doug.SetAndGetRpcCall() wrong")
		return
	}
	// wait test process quit
	time.Sleep(4 * time.Second)
}
