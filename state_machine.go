package raft

import (
	"context"
	"fmt"
	"log"
	"raft/pb"
	"sync"
	"time"
)

type AppendReplyCtx struct {
	peerId int
	err error
	reply *pb.AppendRsp
}

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
	// init leader
	r.nextIndex = make([]int64, len(r.peers))
	r.matchIndex = make([]int64, len(r.peers))
	for i, _ := range r.peers {
		r.nextIndex[i] = r.log.getLastLogIndex() + 1 // noop heartbeat
		r.matchIndex[i] = 0
	}

	hbTmoutC := time.After(time.Millisecond * HbTicks)
	rspC := make(chan *AppendReplyCtx, 128)
	for r.role == Leader {
		select {
		case <- r.stopC:
			log.Printf("node %v stopped", r.me)
			return
		case <- hbTmoutC:
			hbTmoutC = time.After(time.Millisecond * HbTicks)
			go r.doHeartBeat(rspC)
		case replyCtx := <- rspC:
			// handle append response
			r.handleAppendResponse(replyCtx)
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
	electTmoutC := time.After(time.Millisecond * ElectTicks)
	for r.role == Follower {
		select {
		case <- r.stopC:
			log.Printf("node %v stopped", r.me)
			return
		case <- electTmoutC:
			// change to Fellower and reset votedFor
			r.setRole(Candidate)
		case task := <- r.blockQ:
			electTmoutC = time.After(time.Millisecond * ElectTicks)
			// handle
			switch task.req.(type){
			case pb.VoteReq:
				r.handleVote(task)
			case pb.AppendReq:
				r.handleAppend(task)
			}
		}
	}
}

func (r *Raft) candidateLoop(){
	firstElected := false
	electTmoutC := time.After(time.Millisecond * ElectTicks)
	for r.role == Candidate {
		if ! firstElected {
			log.Println("Do first election after become Candidate!")
			firstElected = true
			electTmoutC = time.After(time.Millisecond * ElectTicks)
			go r.doVote()
		}
		select {
		case <- r.stopC:
			log.Printf("node %v stopped", r.me)
			return
		case <- electTmoutC:
			// doVote
			electTmoutC = time.After(time.Millisecond * ElectTicks)
			go r.doVote()
		case task := <- r.blockQ:
			// handle
			switch task.req.(type){
			case pb.VoteReq:
				r.handleVote(task)
			case pb.AppendReq:
				r.handleAppend(task)
			}
		}
	}
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
		r.setCurrentTerm(req.GetTerm())
		r.setRole(Follower)

		// 该状态要不要直接将votedFor更新为candidate呢
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

func (r *Raft) doHeartBeat(rspC chan *AppendReplyCtx){
	log.Println("do HeartBeat now")

	var wg sync.WaitGroup

	for i, _ := range r.peers {
		if i == int(r.me) {
			continue
		}

		go func(idx int) {
			wg.Add(1)
			defer wg.Done()

			es := r.log.log.E[r.nextIndex[i] : ]
			entrys := & pb.Entries{
				E: es,
			}

			req := &pb.AppendReq{
				Term:r.currentTerm,
				LeaderId:r.me,
				PrevLogTerm:r.log.getLastLogTerm(),
				PrevLogIdx:r.log.getLastLogIndex(),
				LeaderCommitLogIdx:r.commitIndex,
				Es:entrys,
			}
			rsp, err := r.peers[idx].Cli.Append(context.Background(), req)
			if err != nil {
				fmt.Println("rcv rsp fail", err)
			}

			replyCtx := & AppendReplyCtx{
				peerId:idx,
				err:err,
				reply:rsp,
			}
			rspC <- replyCtx
		}(i)
	}

	wg.Wait()
}

func (r *Raft)handleAppendResponse(replyCtx *AppendReplyCtx){
	if replyCtx.err != nil {
		log.Println("Receive a err response!")
		return
	}

	if replyCtx.reply.RetCode == RetCodeOk {
		r.matchIndex[replyCtx.peerId] = r.nextIndex[replyCtx.peerId]
		r.nextIndex[replyCtx.peerId] += 1

		// whether commitIndex can be increased
		r.doCommit()
	} else {
		log.Printf("Reduce peer %v nextIndex %v", replyCtx.peerId,
			r.nextIndex[replyCtx.peerId])
		if r.nextIndex[replyCtx.peerId] > 0 {
			r.nextIndex[replyCtx.peerId] -= 1
		}
	}
}


func (r *Raft) doVote(){
	// 1. increase currentTerm
	r.currentTerm = r.currentTerm + 1

	// 2. vote to self
	r.log.meta.VotedFor = r.me
	r.log.PersistMeta(r.log.meta)

	// 3. send vote request
	sucNum := 1
	var wg sync.WaitGroup
	for i, _ := range r.peers {
		if i == int(r.me) {
			continue
		}

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			req := &pb.VoteReq{
				Term:r.currentTerm,
				CandiId:r.me,
				LastLogTerm:r.log.getLastLogTerm(),
				LastLogIdx:r.log.getLastLogIndex(),
			}

			rsp, err := r.peers[idx].Cli.Vote(context.Background(), req)
			if err != nil {
				fmt.Println("rcv rsp fail", err)
				return
			}

			if rsp.GetVoteGranted() == RetCodeOk {
				sucNum += 1
			}else{
				log.Println("Not get vote")
			}
		}(i)
	}

	wg.Wait()
	if sucNum >= r.quorum() {
		r.setRole(Leader)
	}
}