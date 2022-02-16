package main

import serverPKG "code/src/server"

func main() {
	server := serverPKG.Server{NameGroup: make(map[string]string), NameIps: make(map[string]string),
		GroupMessage: make(map[string][]string), NameRead: make(map[string]int),
	}
	server.Proceed()
}
