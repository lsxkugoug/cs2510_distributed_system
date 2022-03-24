package rpc

// OperationArgs ClientSendMessageArgs
// When a client send a message to the group, the server should store it
type OperationArgs struct {
	Operation int  // 0 means update or new, 1 means retrieve data
	IsClient  bool // whether the msg is come from client, if yes, after update, broadcast, if no, not
	Key       string
	Value     string
	Vector    []int // if client send, Vector == nil
}

type OperationReply struct {
	Operation  int // 0 means update or new, 1 means retrieve data
	DataVector [][]int
	DataValue  []string
}
