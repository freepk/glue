package main

import (
	"testing"
)

func intsEq(a, b []int) bool {
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func bytesEq(a, b []byte) bool {
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseUint(t *testing.T) {
	if _, n := parseUint(nil); n != 0 {
		t.Fail()
	}
	if _, n := parseUint([]byte{}); n != 0 {
		t.Fail()
	}
	if _, n := parseUint([]byte("a")); n != 0 {
		t.Fail()
	}
	// TODO: max int overload
	//if v, n = ParseUint([]byte("123456789")); n != 0 {
	//	t.Fail()
	//}
	if v, n := parseUint([]byte("1234")); (n != 4) || (v != 1234) {
		t.Fail()
	}
	if v, n := parseUint([]byte("0123a")); (n != 4) || (v != 123) {
		t.Fail()
	}
	if v, n := parseUint([]byte("1234a")); (n != 4) || (v != 1234) {
		t.Fail()
	}
}

func TestAppendUint(t *testing.T) {
	if !bytesEq(appendUint(nil, 0), []byte("0")) {
		t.Fail()
	}
	if !bytesEq(appendUint(nil, 10), []byte("10")) {
		t.Fail()
	}
	if !bytesEq(appendUint(nil, 1234), []byte("1234")) {
		t.Fail()
	}
}

func TestParseUints(t *testing.T) {
	if _, n := parseUints(nil, nil); n != 0 {
		t.Fail()
	}
	if _, n := parseUints(nil, []byte{}); n != 0 {
		t.Fail()
	}
	if _, n := parseUints(nil, []byte("a")); n != 0 {
		t.Fail()
	}
	if _, n := parseUints(nil, []byte(";;;")); n != 3 {
		t.Fail()
	}
	if r, n := parseUints(nil, []byte("1234")); n != 4 || !intsEq(r, []int{1234}) {
		t.Fail()
	}
	if r, n := parseUints(nil, []byte("1234;")); n != 5 || !intsEq(r, []int{1234}) {
		t.Fail()
	}
	if r, n := parseUints(nil, []byte("1234;5678")); n != 9 || !intsEq(r, []int{1234, 5678}) {
		t.Fail()
	}
}
