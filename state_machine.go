package raft

import (
	"log"
	"raft/pb"
)

func (r *Raft) eventLoop(){
	switch r.role {
	case Leader:
		r.leaderLoop()
	case Follower:
		r.followerLoop()
	case Candidate:
		r.candidateLoop()
	}
}

func (r *Raft) leaderLoop(){
	for r.role == Leader {
		select {
		case <- r.stopC:
			log.Printf("node %v stopped", r.me)
			return
		case task := <- r.blockQ:
			// handle
			switch task.req.(type){
			case pb.VoteReq:
				r.handleVote(task)
			case pb.AppendReq:
				r.handleAppend(task)
			//case
			}
		}
	}
}

func (r *Raft) followerLoop(){

}

func (r *Raft) candidateLoop(){

}

func (r *Raft) handleVote(task message) {
	req := task.req.(pb.VoteReq)
	rsp := task.rsp.(pb.VoteRsp)

	defer func() {
		task.errC <- nil
	}()

	if req.GetTerm() < r.currentTerm {
		log.Printf("Local currentTerm:%d, req Term:%v", r.currentTerm, req.GetTerm())
		rsp.Term = r.currentTerm
		rsp.VoteGranted = RetCodeErr
		return
	}

	if req.GetTerm() > r.currentTerm {
		log.Printf("Update term from %d to %d", r.currentTerm, req.GetTerm())
		r.currentTerm = req.GetTerm()
	}

	if r.log.getVotedFor() != EmptyVotedFor && r.log.getVotedFor() != req.GetCandiId() {
		log.Printf("Local votedFor:%d, req candidate:%v", r.log.getVotedFor(),
			req.GetCandiId())
		rsp.Term = r.currentTerm
		rsp.VoteGranted = RetCodeErr
		return
	}

	if req.GetLastLogTerm() > r.log.getLastLogTerm() ||
		req.GetLastLogTerm() == r.log.getLastLogTerm() && req.GetLastLogIdx() >= r.log.getLastLogIndex() {
		meta := & pb.Meta{VotedFor:req.GetCandiId()}
		r.log.PersistMeta(meta)

		rsp.VoteGranted = RetCodeOk
		rsp.Term = r.currentTerm
		return
	}

	rsp.VoteGranted = RetCodeErr
	rsp.Term = r.currentTerm
	return
}

func (r *Raft) handleAppend(task message) {
	req := task.req.(pb.AppendReq)
	rsp := task.rsp.(pb.AppendRsp)

	defer func() {
		task.errC <- nil
	}()

	if req.GetTerm() < r.currentTerm {
		log.Printf("Local currentTerm:%d, req Term:%v", r.currentTerm, req.GetTerm())
		rsp.Term = r.currentTerm
		rsp.RetCode = RetCodeErr
		return
	}
	
	if r.log.log.E[req.GetPrevLogIdx()].GetTerm() != req.GetPrevLogTerm() {
		log.Printf("Local currentTerm:%d, req Term:%v", r.currentTerm, req.GetTerm())
		rsp.Term = r.currentTerm
		rsp.RetCode = RetCodeErr
		return
	}
}