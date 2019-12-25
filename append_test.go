package raft

import "testing"

func TestArr(t *testing.T){
	arr := []int{1, 2, 3, 4}
	a1 := arr[1:3]
	for _, v := range a1 {
		t.Logf("v = %v", v)
	}

	a1Len := len(a1)
	a2 := a1[a1Len:]
	t.Logf("a2Len = %v", len(a2))
}
