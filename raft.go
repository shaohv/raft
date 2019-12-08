package raft

const (
	RetCodeOk  = 0
	RetCodeErr = -1
)

const (
	DefaultTerm = -1
	DefaultLeader = -1
	EmptyVotedFor = -1
)

const (
	Leader = 0
	Follower = 1
	Candidate = 2
)

type message struct {
	// receive all type requestes
	req interface{}

	// the caller
	rsp interface{}
	errC chan error
}

type Raft struct {
	// permently
	currentTerm int64
	//votedFor int32
	//log[]
	log *Log

	commitIndex int64
	lastApplied int64

	nextIndex []int64
	matchIndex []int64

	role int8
	me int16
	stopC chan bool
	blockQ chan message
}

func NewRaft() *Raft{
	return &Raft{
		//currentTerm:DefaultTerm,
		//votedFor:DefaultLeader,
		currentTerm:0,
		log : NewLog(),
		commitIndex:-1,
		lastApplied:-1,
		stopC:make(chan bool, 1),
		blockQ:make(chan message, 128),
	}
}

func (r *Raft) deliver(req interface{}, rsp interface{}, errC chan error) error{
	task := message{req:req, rsp:rsp, errC:errC}
	r.blockQ <- task

	select {
	case err := <- errC:
		return err
	}
}