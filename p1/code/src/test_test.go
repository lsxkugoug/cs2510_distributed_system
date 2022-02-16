package main

import (
	serverPackage "code/src/server"
	"os/exec"
	"testing"
	"time"
)

// In this test Alice online first, and send the "I am Alice"
// Then Bob online, he would receive the unread message: "I am Alice", and reply "nice to meet you"
// Alice receives the message
func Test1(t *testing.T) {
	// 1.run server
	server := serverPackage.Server{NameGroup: make(map[string]string), NameIps: make(map[string]string),
		GroupMessage: make(map[string][]string), NameRead: make(map[string]int),
	}

	go server.Proceed()
	time.Sleep(1 * time.Second)

	// 2.run Alice, Bob, Chad
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test1/AliceMain.go").Start()
	time.Sleep(5 * time.Second)
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test1/BobMain.go").Start()
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test1/ChadMain.go").Start()

	// waite bob and chad write the log
	time.Sleep(5 * time.Second)

}

//Alice, Bob, and Chad are online.
//Bob sends a message to all, Chad and Alice receives the message.
//Alice sends a message to all, Bob and Chad receives it (but not Alice).
//Doug, not part of the group, joins the server but receives no message.
func Test2(t *testing.T) {
	// 1.run server
	server := serverPackage.Server{NameGroup: make(map[string]string), NameIps: make(map[string]string),
		GroupMessage: make(map[string][]string), NameRead: make(map[string]int),
	}
	go server.Proceed()
	time.Sleep(1 * time.Second)

	exec.Command("/bin/bash", "-c", "go run ./testEntry/test2/AliceMain.go").Start()
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test2/BobMain.go").Start()
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test2/ChadMain.go").Start()
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test2/DougMain.go").Start()

	// waite bob and chad write the log
	time.Sleep(5 * time.Second)

}

//Three member Alice, Bob, Chad. All of them in  the same group
//Alice send a message "I am Alice, nice to meet you" to the group, Bob would receive it.
//Then Bob offline.
//Chad  online and send the message "I am Chad, nice to meet you"to the group
//Then Bob online, he only would receive the "I am Chad, nice to meet you", because  "I am Alice, nice to meet you" he have already read.
func Test3(t *testing.T) {
	// 1.run server
	server := serverPackage.Server{NameGroup: make(map[string]string), NameIps: make(map[string]string),
		GroupMessage: make(map[string][]string), NameRead: make(map[string]int),
	}
	go server.Proceed()
	time.Sleep(1 * time.Second)

	exec.Command("/bin/bash", "-c", "go run ./testEntry/test3/AliceMain.go").Start()
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test3/BobMain.go").Start()
	// Bob would offline
	time.Sleep(5 * time.Second)
	// Chad send message, Bob only read the Chad's history message, but not Alice,
	//since Alice's message is already read by Bob
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test3/ChadMain.go").Start()
	time.Sleep(1 * time.Second)
	exec.Command("/bin/bash", "-c", "go run ./testEntry/test3/BobMain.go").Start()

	// waite bob and chad write the log
	time.Sleep(5 * time.Second)

}
