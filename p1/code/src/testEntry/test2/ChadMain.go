package main

import (
	clientPackage "code/src/client"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("./testEntry/test2/log/Chad.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	Chad := clientPackage.Client{UserName: "Chad", UserGroup: "1",
		ServerIpPort: clientPackage.GetOutboundIP("1234"), ClientPort: "1222",
		ClientIpPort: clientPackage.GetOutboundIP("1222")}
	conn, err, listener := Chad.SetAndGetConn()
	defer listener.Close()
	defer conn.Close()
	if err != nil {
		fmt.Println("Chad.SetAndGetRpcCall() wrong")
		return
	}
	// wait test process quit
	time.Sleep(4 * time.Second)

}
