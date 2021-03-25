package main

import (
	"github.com/valyala/fasthttp"
	"sync"
)

type worker struct {
	pool *workerPool
	wait *sync.WaitGroup
	rq   []*fasthttp.Request
	rs   []*fasthttp.Response
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.wait = new(sync.WaitGroup)
	p.workers <- w
	return w
}

func (w *worker) reset() {
	w.rq = w.rq[:0]
	w.rs = w.rs[:0]
}

func doAsync(wait *sync.WaitGroup, rq *fasthttp.Request, rs *fasthttp.Response) {
	defer wait.Done()
	fasthttp.Do(rq, rs)
}

func (w *worker) run(buf []byte, urls []string) []byte {
	n := len(urls)
	w.wait.Add(n)
	for i := 0; i < n; i++ {
		rs := fasthttp.AcquireResponse()
		rq := fasthttp.AcquireRequest()
		w.rq = append(w.rq, rq)
		w.rs = append(w.rs, rs)
		rq.SetRequestURI(urls[i])
		go doAsync(w.wait, rq, rs)
	}
	w.wait.Wait()
	for i := 0; i < n; i++ {
		buf = append(buf, w.rs[i].Body()...)
		fasthttp.ReleaseRequest(w.rq[i])
		fasthttp.ReleaseResponse(w.rs[i])
	}
	return buf
}

func (w *worker) release() {
	w.pool.workers <- w
}

type workerPool struct {
	workers chan *worker
}

func newWorkerPool(n int) *workerPool {
	p := new(workerPool)
	p.workers = make(chan *worker, n)
	for i := 0; i < n; i++ {
		newWorker(p)
	}
	return p
}

func (p *workerPool) acquire() *worker {
	w := <-p.workers
	w.reset()
	return w
}
