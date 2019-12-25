package raft

import (
	"bufio"
	"errors"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	RetCodeOk  = 0
	RetCodeErr = -1
)

const (
	DefaultTerm = 0
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

var RoleStr = [3]string{"Leader", "Follower", "Candidate"}

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
	me int32
	grpcServer *grpc.Server
	stopC chan bool
	blockQ chan message
}

func (r *Raft) setCurrentTerm(t int64){
	r.currentTerm = t
}

func (r *Raft) setRole(ro int8){
	older := "None"
	if r.role >= 0 {
		older = RoleStr[r.role]
	}

	newer := "None"
	if ro >= 0 {
		newer = RoleStr[ro]
	}
	log.Println("Role change:%s -> %s", older, newer)

	r.role = ro
}

func (r *Raft) quorum() int {
	return len(r.peers) / 2 + 1;
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

/*
	conf 配置文件路径
cluster=ip1:port1;ip2:port2;ip3:port3
me=nid
*/
func (r *Raft) Startup(conf string) error{
	if err := r.readConf(conf); err != nil {
		log.Println("Read config failed, %v", err)
		return err
	}

	addr := r.peers[r.me].Ip + r.peers[r.me].Port
	serverC := make(chan *grpc.Server, 1)
	StartGrpcServer(addr, serverC)
	r.grpcServer = <- serverC

	for i, node := range r.peers {
		for {
			addr = node.Ip + node.Port
			p, err := NewPeer(addr)
			if err != nil {
				log.Println("Maybe peer node not startup %v", err)
				continue
			}
			r.peers[i].Cli = p.Cli
			r.peers[i].Conn = p.Conn
			break
		}
	}

	return nil
}

func (r *Raft) readConf(conf string) error {
	f, err := os.Open(conf)
	if err != nil {
		log.Println("Open conf file %v failed, %v", conf, err)
		return errors.New("Open config file failed!")
	}
	defer f.Close()

	bufR := bufio.NewReader(f)
	cluster, err := bufR.ReadString('\n')
	if err != nil {
		log.Println("Read cluster info failed, %v", err)
		return errors.New("Read cluster information failed!")
	}
	log.Println("Cluster information:%v", cluster)

	me, err := bufR.ReadString('\n')
	if err != nil {
		log.Println("Read me failed, %v", err)
		return errors.New("Read me failed!")
	}
	log.Println("My nid is %v", me)

	nodes := strings.Split(cluster, ",")
	for _,v := range nodes {
		log.Println(v)
		ipport := strings.Split(v,":")
		node := PeerCli{
			Ip:ipport[0],
			Port:ipport[1],
		}
		r.peers = append(r.peers, node)
	}

	nid, err := strconv.Atoi(me)
	if err != nil {
		log.Println("Atoi failed %v", err)
		return errors.New("Atoi failed")
	}

	r.me = int32(nid)

	return nil
}

func (r *Raft) deliver(req interface{}, rsp interface{}, errC chan error) error{
	task := message{req:req, rsp:rsp, errC:errC}
	r.blockQ <- task

	select {
	case err := <- errC:
		return err
	}
}

func (r *Raft) doCommit(){
	mini := r.log.getLastLogIndex()

	for j := mini; j > 0; j-- {
		cnt := 1
		for i, v := range r.matchIndex {
			if i != int(r.me) && v == j {
				cnt ++
			}
		}

		if cnt >= r.quorum() {
			log.Println("Update commitIndex %v -> %v", r.commitIndex, mini)
			r.commitIndex = mini
			break
		}
	}
}