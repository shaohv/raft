package raft

const (
	RetCodeOk  = 0
	RetCodeErr = -1
)

const (
	DefaultTerm = 0
	DefaultLeader = -1
	EmptyVotedFor = -1
	DefaultLastApplied = 0
	DefaultCommitIndex = 0
)

const (
	Leader = 0
	Follower = 1
	Candidate = 2
)

const (
	HbTicks = 3
	ElectTicks = 150
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

	commitIndex int64 // initialized to 0, so the []log first index is 1
	lastApplied int64

	nextIndex []int64
	matchIndex []int64

	role int8
	peers []PeerCli
	me int16
	stopC chan bool
	blockQ chan message
}

func (r *Raft) setCurrentTerm(t int64){
	r.currentTerm = t
}

func (r *Raft) setRole(ro int8){
	r.role = ro
}

func NewRaft() *Raft{
	r := &Raft{
		//currentTerm:DefaultTerm,
		//votedFor:DefaultLeader,
		currentTerm:DefaultTerm,
		log : NewLog(),
		commitIndex:DefaultCommitIndex,
		lastApplied:DefaultLastApplied,
		role: Follower,
		stopC:make(chan bool, 1),
		blockQ:make(chan message, 128),
	}

	r.log.LoadMeta()
	r.log.LoadLog()

	r.currentTerm = r.log.getLastLogTerm()

	return r
}

func (r *Raft) deliver(req interface{}, rsp interface{}, errC chan error) error{
	task := message{req:req, rsp:rsp, errC:errC}
	r.blockQ <- task

	select {
	case err := <- errC:
		return err
	}
}