package pb

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"github.com/golang/protobuf/proto"
)


func TestLogPb(t *testing.T){
	log.Println("entry TestLogPb ...")
	es := &Entries{}

	e1 := & Entry{Term:1,LogIdx:1,Key:"1111",Value:"1111"}
	e2 := & Entry{Term:2,LogIdx:2,Key:"2222",Value:"2222"}

	es.E = append(es.E, e1)
	es.E = append(es.E, e2)

	data, err :=proto.Marshal(es)
	if err != nil {
		t.Error("marshal failed")
		return
	}

	err = ioutil.WriteFile("./TestLogPb", data, os.ModePerm)
	if nil != err{
		t.Error("write failed")
		return
	}

	rdata, err := ioutil.ReadFile("./TestLogPb")
	rEs := & Entries{}
	err = proto.Unmarshal(rdata, rEs)
	if err != nil {
		t.Error("read data failed")
		return
	}

	if len(rEs.E) != 2 {
		t.Errorf("real len:%v", len(rEs.E))
		return
	}

	if rEs.E[0].Term != 1 {
		t.Errorf("real term:%v", rEs.E[0].Term)
		return
	}

	if rEs.E[1].Term != 2 {
		t.Errorf("real term:%v", rEs.E[1].Term)
		return
	}
	log.Println("leave TestLogPb ...")
}