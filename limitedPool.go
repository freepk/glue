package main

type limitedPool struct {
	queue chan interface{}
}

func newLimitedPool(newFunc func() interface{}, limit int) *limitedPool {
	p := new(limitedPool)
	p.queue = make(chan interface{}, limit)
	for i := 0; i < limit; i++ {
		p.queue <- newFunc()
	}
	return p
}

func (p *limitedPool) get() interface{} {
	return <-p.queue
}

func (p *limitedPool) put(elem interface{}) {
	p.queue <- elem
}
