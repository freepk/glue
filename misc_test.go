package main

import "testing"

func TestDedupInts(t *testing.T) {
	if !intsEq(dedupInts([]int{}), []int{}) {
		t.Fail()
	}
	if !intsEq(dedupInts([]int{1, 1}), []int{1}) {
		t.Fail()
	}
	if !intsEq(dedupInts([]int{1, 2, 2, 3, 4, 4, 4}), []int{1, 2, 3, 4}) {
		t.Log(dedupInts([]int{1, 2, 2, 3, 4, 4, 4}))
		t.Fail()
	}
}
