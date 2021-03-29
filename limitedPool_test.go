package main

import "testing"

func TestNewLimitedPool(t *testing.T) {
	p := newLimitedPool(func() interface{} {
		return nil
	}, 4)

	w0 := p.get()
	w1 := p.get()
	w2 := p.get()
	w3 := p.get()

	go func() { w4 := p.get(); p.put(w4) }()

	// release for w4
	p.put(w0)

	// wait for w4.release()
	w0 = p.get()
	p.put(w0)

	p.put(w1)
	p.put(w2)
	p.put(w3)
}
