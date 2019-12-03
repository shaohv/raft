package raft

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
	votedFor int32
	//log[]

	commitIndex int64
	lastApplied int64

	nextIndex []int64
	matchIndex []int64

	blockQ chan message
}

func (r *Raft) deliver(req interface{}, rsp interface{}, errC chan error) error{
	task := message{req:req, rsp:rsp, errC:errC}
	r.blockQ <- task

	select {
	case err := <- errC:
		return err
	}
}