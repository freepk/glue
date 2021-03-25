package main

import (
	"testing"
)

func TestWorkerPool(t *testing.T) {
	p := newWorkerPool(4)
	w0 := p.acqure()
	w1 := p.acqure()
	w2 := p.acqure()
	w3 := p.acqure()

	go func() { w4 := p.acqure(); w4.release() }()

	// release for w4
	w0.release()

	// wait for w4.release()
	w0 = p.acqure()
	w0.release()
	w1.release()
	w2.release()
	w3.release()
}
