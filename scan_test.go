package main

import (
	"testing"
)

func isEqual(a, b []int) bool {
	size := len(a)
	if size != len(b) {
		return false
	}
	for i := 0; i < size; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestScanNum(t *testing.T) {
	v := 0
	n := 0
	if v, n = scanNum(nil); n != 0 {
		t.Fail()
	}
	if v, n = scanNum([]byte{}); n != 0 {
		t.Fail()
	}
	if v, n = scanNum([]byte("a")); n != 0 {
		t.Fail()
	}
	// TODO: max int overload
	//if v, n = scanNum([]byte("123456789")); n != 0 {
	//	t.Fail()
	//}
	if v, n = scanNum([]byte("1234")); (n != 4) || (v != 1234) {
		t.Fail()
	}
	if v, n = scanNum([]byte("0123a")); (n != 4) || (v != 123) {
		t.Fail()
	}
	if v, n = scanNum([]byte("1234a")); (n != 4) || (v != 1234) {
		t.Fail()
	}
}

func TestScanNums(t *testing.T) {
	r := make([]int, 0, 1024)
	x := r
	n := 0
	if x, n = scanNums(r, nil); n != 0 {
		t.Fail()
	}
	if x, n = scanNums(r, []byte{}); n != 0 {
		t.Fail()
	}
	if x, n = scanNums(r, []byte("a")); n != 0 {
		t.Fail()
	}
	if x, n = scanNums(r, []byte(";;;")); n != 3 {
		t.Fail()
	}
	if x, n = scanNums(r, []byte("1234")); n != 4 || !isEqual(x, []int{1234}) {
		t.Fail()
	}
	if x, n = scanNums(r, []byte("1234;")); n != 5 || !isEqual(x, []int{1234}) {
		t.Fail()
	}
	if x, n = scanNums(r, []byte("1234;5678")); n != 9 || !isEqual(x, []int{1234, 5678}) {
		t.Fail()
	}
}

func TestDedupNums(t *testing.T) {
	if !isEqual(dedupNums([]int{}), []int{}) {
		t.Fail()
	}
	if !isEqual(dedupNums([]int{1, 1}), []int{1}) {
		t.Fail()
	}
	if !isEqual(dedupNums([]int{1, 2, 2, 3, 4, 4, 4}), []int{1, 2, 3, 4}) {
		t.Fail()
	}
}
