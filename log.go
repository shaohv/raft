package raft

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"raft/pb"
)

const (
	PaddingKey = "paddingKey"
	PaddingValue = "paddingValue"
)
type Log struct {
	meta *pb.Meta
	log *pb.Entries

	metaFile string
	logFile string
}

func NewLog() *Log {
	return & Log{
		meta : & pb.Meta{VotedFor:EmptyVotedFor},
		log : & pb.Entries{},
		metaFile:"meta.dat",
		logFile:"log.dat",
	}
}

func (l *Log) LoadMeta() {
	f, err := os.Open(l.metaFile)

	if err != nil && os.IsNotExist(err) {
		// create && write initiatie meta
		l.PersistMeta(l.meta)
	}else if err == nil {
		// read
		defer f.Close()

		data, err := ioutil.ReadFile(l.metaFile)
		err = proto.Unmarshal(data, l.meta)
		if err != nil {
			panic("read meta failed")
		}
	}else{
		panic("open " + l.metaFile + " failed!")
	}
}

func (l *Log) LoadLog(){
	f, err := os.Open(l.logFile)
	if err != nil && os.IsNotExist(err) {
		// create && write initiatie log
		paddingEntry := & pb.Entry{
			Term:DefaultTerm,
			LogIdx:0,
			Key:"paddingKey",
			Value:"paddingValue",
		}
		l.log.E = append(l.log.E, paddingEntry)
		l.PersistLog(l.log)
	}else if err == nil {
		// read
		defer f.Close()

		data, err := ioutil.ReadFile(l.logFile)
		err = proto.Unmarshal(data, l.log)
		if err != nil {
			panic("read meta failed")
		}
	}else{
		panic("open " + l.metaFile + " failed!")
	}
}

func (l *Log) getLastLogTerm() int64 {
	//return l.currentTerm
	size := len(l.log.E)
	if size == 0 {
		return 0
	} else {
		return l.log.E[size - 1].Term
	}
}

func (l *Log) getLastLogIndex() int64 {
	//return l.currentTerm
	size := len(l.log.E)
	if size == 0 {
		return 0
	} else {
		return l.log.E[size - 1].LogIdx
	}
}

func (l *Log) getVotedFor() int32 {
	//return l.votedFor
	return l.meta.VotedFor
}

func (l *Log) PersistLog(data *pb.Entries) {
	d, err := proto.Marshal(data)
	if err != nil {
		panic("marshal data")
	}

	err = ioutil.WriteFile(l.logFile, d, os.ModePerm)
	if err != nil {
		panic("write file")
	}

	l.log = data
}

func (l *Log) PersistMeta(meta *pb.Meta) {
	d, err := proto.Marshal(meta)
	if err != nil {
		panic("marshal data")
	}

	err = ioutil.WriteFile(l.metaFile, d, os.ModePerm)
	if err != nil {
		panic("write file")
	}

	l.meta = meta
}