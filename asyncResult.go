package main

type asyncResult struct {
	ch chan []byte
	r  [][]byte
	fn func([]byte, []byte) []byte
}

func newAsyncResult(fn func([]byte, []byte) []byte) *asyncResult {
	ar := new(asyncResult)
	ar.ch = make(chan []byte)
	ar.r = make([][]byte, 256)
	ar.fn = fn
	return ar
}

func exec(ch chan []byte, a, b []byte, fn func([]byte, []byte) []byte) {
	ch <- fn(a, b)
}

func (ar *asyncResult) exec(a []byte, b [][]byte) []byte {
	for i := 0; i < len(b); i++ {
		go exec(ar.ch, ar.r[i][:0], b[i], ar.fn)
	}
	for i := 0; i < len(b); i++ {
		ar.r[i] = <-ar.ch
		a = append(a, ar.r[i]...)
	}
	return a
}
