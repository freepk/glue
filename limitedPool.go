package main

type limitedPool struct {
	ch chan interface{}
}

func newLimitedPool(fn func() interface{}, limit int) *limitedPool {
	p := new(limitedPool)
	p.ch = make(chan interface{}, limit)
	for i := 0; i < limit; i++ {
		p.ch <- fn()
	}
	return p
}

func (p *limitedPool) get() interface{} {
	return <-p.ch
}

func (p *limitedPool) put(i interface{}) {
	p.ch <- i
}
