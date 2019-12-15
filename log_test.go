package raft

import (
	"os"
	"raft/pb"
	"testing"
)

func TestReadMeta(t *testing.T){
	log := NewLog()

	defer func() {
		err := os.Remove(log.metaFile)
		if err != nil {
			t.Errorf("Remove %s failed!, %v", log.metaFile, err)
		}
	}()

	log.LoadMeta()
	if log.meta.VotedFor != EmptyVotedFor {
		t.Errorf("votedFor:%v, %v", log.meta.VotedFor, EmptyVotedFor)
		return
	}

	t.Log("Successful!")
}

func TestReadMetaWithUpd(t *testing.T){
	log := NewLog()

	defer func() {
		err := os.Remove(log.metaFile)
		if err != nil {
			t.Errorf("Remove %s failed!, %v", log.metaFile, err)
		}
	}()

	log.LoadMeta()
	if log.meta.VotedFor != EmptyVotedFor {
		t.Errorf("votedFor:%v, %v", log.meta.VotedFor, EmptyVotedFor)
		return
	}

	log.meta.VotedFor = 10
	log.PersistMeta(log.meta)
	log.meta.VotedFor = EmptyVotedFor
	log.LoadMeta()
	if log.meta.VotedFor != 10 {
		t.Errorf("votedFor:%v, %v", log.meta.VotedFor, EmptyVotedFor)
		return
	}

	t.Log("Successful!")
}

func TestReadLog(t *testing.T){
	log := NewLog()

	defer func() {
		err := os.Remove(log.logFile)
		if err != nil {
			t.Errorf("Remove %s failed!, %v", log.logFile, err)
		}
	}()

	log.LoadLog()
	if len(log.log.E) != 1 {
		t.Errorf("log length:%v", len(log.log.E))
		return
	}

	padding := log.log.E[0]
	if padding.Term != DefaultTerm || padding.LogIdx != 0 {
		t.Errorf("Term:%v, LogIndex:%v", padding.Term, padding.LogIdx)
		return
	}

	if padding.Key != PaddingKey || padding.Value != PaddingValue {
		t.Errorf("padding key :%v, padding value : %v", padding.Key, padding.Value)
		return
	}
	t.Logf("padding Term:%v, LogIndex:%v, Key:%v, Value:%v",
		padding.Term, padding.LogIdx, padding.Key, padding.Value)
	t.Log("Successful!")
}

func TestReadLogWithAppend(t *testing.T){
	log := NewLog()

	defer func() {
		err := os.Remove(log.logFile)
		if err != nil {
			t.Errorf("Remove %s failed!, %v", log.logFile, err)
		}
	}()

	log.LoadLog()
	if len(log.log.E) != 1 {
		t.Errorf("log length:%v", len(log.log.E))
		return
	}

	e := & pb.Entry{
		Term:2,
		LogIdx:100,
		Key:"key1",
		Value:"value1",
	}

	e2 := & pb.Entry{
		Term:3,
		LogIdx:1001,
		Key:"key2",
		Value:"value2",
	}
	log.log.E = append(log.log.E, e, e2)
	log.PersistLog(log.log)
	log.log = & pb.Entries{}
	log.LoadLog()

	if len(log.log.E) != 3 {
		t.Errorf("log length:%v", len(log.log.E))
		return
	}

	a := log.log.E[1]
	b := log.log.E[2]

	if a.Term != 2 || a.LogIdx != 100 || b.Term != 3 || b.LogIdx != 1001 {
		t.Logf("a Term:%v, LogIndex:%v, b Term:%v, LogIndex:%v",
			a.Term, a.LogIdx, b.Term, b.LogIdx)
		return
	}

	t.Logf("a Term:%v, LogIndex:%v, Key:%v, Value:%v",
		a.Term, a.LogIdx, a.Key, a.Value)
	t.Logf("b Term:%v, LogIndex:%v, Key:%v, Value:%v",
		b.Term, b.LogIdx, b.Key, b.Value)

	if log.getLastLogTerm() != 3 || log.getLastLogIndex() != 1001 {
		t.Errorf("LastLogTerm:%v, LastLogIndex:%v", log.getLastLogTerm(), log.getLastLogIndex())
		return
	}
	t.Log("Successful!")
}