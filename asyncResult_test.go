package main

import "testing"

func test(a, b []byte) []byte {
	return append(a, b...)
}

func TestAsyncResult(t *testing.T) {
	ar := newAsyncResult(test)
	one := make([]byte, 0, 256)
	two := [][]byte{[]byte(`{1}`), []byte(`{2}`), []byte(`{3}`)}
	for i := 0; i < 1024; i++ {
		one = ar.exec(one[:0], two)
		switch {
		case bytesEq(one, []byte(`{1}{2}{3}`)):
		case bytesEq(one, []byte(`{1}{3}{2}`)):
		case bytesEq(one, []byte(`{2}{1}{3}`)):
		case bytesEq(one, []byte(`{2}{3}{1}`)):
		case bytesEq(one, []byte(`{3}{1}{2}`)):
		case bytesEq(one, []byte(`{3}{2}{1}`)):
		default:
			t.Fail()
		}
	}
}

func BenchmarkAsyncResult(b *testing.B) {
	ar := newAsyncResult(test)
	one := make([]byte, 0, 256)
	two := [][]byte{[]byte(`{1}`), []byte(`{2}`), []byte(`{3}`)}
	for i := 0; i < b.N; i++ {
		one = ar.exec(one[:0], two)
	}
}
