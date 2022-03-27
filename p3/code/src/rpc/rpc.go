package rpc

// OperationArgs ClientSendMessageArgs
// When a client send a message to the group, the server should store it
type OperationArgs struct {
	Operation int // 0 means update or new, 1 means retrieve data
	IsLeader  bool
	Key       string
	Value     string
	IsCommit  bool
	IsRead    bool // client -> replica->other replicas, and setIsRead == true
}

type OperationReply struct {
	Operation int    // 0 means update or new, 1 means retrieve data
	Value     string //
	IsReady   bool
}

type LeaderElectionArgs struct {
	SenderId   int // sender's id
	Type       int // 0 means I want to make election, 1 means I successfully become the leader
	SenderTerm int
}

type LeaderElectionReply struct {
	IsOk bool //only for leader election message, if my id > senderId, send ok = true, else false
}
