package main

import (
	"testing"
)

func TestWorkerPool(t *testing.T) {
	p := newWorkerPool(4)

	w0 := p.acquire()
	w1 := p.acquire()
	w2 := p.acquire()
	w3 := p.acquire()

	go func() { w4 := p.acquire(); w4.release() }()

	// release for w4
	w0.release()

	// wait for w4.release()
	w0 = p.acquire()
	w0.release()
	w1.release()
	w2.release()
	w3.release()
}

func TestWorkerPool128(t *testing.T) {
	newWorkerPool(128)
}
