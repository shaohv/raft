package raft

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"raft/pb"
)

//type PeerServerImpl struct {
//}

func (r *Raft) Append(ctx context.Context, req *pb.AppendReq) (*pb.AppendRsp, error) {
	log.Println("Handle Append Message!")

	rsp := & pb.AppendRsp{}
	err := r.deliver(req, rsp, make(chan error))
	if err != nil {
		rsp.RetCode = RetCodeErr
		return rsp, err
	}
	return rsp, nil
}

func (r *Raft) Vote(ctx context.Context, req *pb.VoteReq) (*pb.VoteRsp, error) {
	log.Println("Handle Vote Message!")
	//rsp := &pb.VoteRsp{Term:2,VoteGranted:-1}
	rsp := &pb.VoteRsp{}

	err := r.deliver(req, rsp, make(chan error, 1))
	if err != nil {
		rsp.VoteGranted = RetCodeErr
		return rsp, err
	}
	return rsp, nil
}

func StartGrpcServer(addr string, servC chan *grpc.Server) {
	log.Println("Start grpc server now")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterPeerServer(s, &Raft{})
	servC <- s
	s.Serve(lis)
}

type PeerCli struct {
	Conn *grpc.ClientConn
	Cli  pb.PeerClient
	Ip string
	Port string
}

// GetPeer: setup connections with grpc server and return grpc client
func NewPeer(addr string) (*PeerCli, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("conn addr:%v fail:%v", addr, err)
		return nil, err
	}

	client := pb.NewPeerClient(conn)
	cli := &PeerCli{
		Conn: conn,
		Cli:  client,
	}
	return cli, nil
}
