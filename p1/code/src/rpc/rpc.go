package rpc

// ServerSendMessageArgs
//When a client send a message to the server,
//the server should send the message to the all clients who are in the same group.
type ServerSendMessageArgs struct {
	Messages []string
}

type ServerSendMessageReply struct {
	Status bool // 0 success, -1 wrong
}

// ClientGetMessageArgs
// When a client online, he needs to receive all the message he hasn't read
type ClientGetMessageArgs struct {
	UserName     string
	UserGroup    string
	ClientIpPort string //used to allow server flooding
}

type ClientGetMessageReply struct {
	UnreadMessages []string
}

// ClientSendMessageArgs
// When a client send a message to the group, the server should store it
type ClientSendMessageArgs struct {
	UserName string
	Message  string
}

type ClientSendMessageReply struct {
	Status bool // 0 success, -1 wrong
}
