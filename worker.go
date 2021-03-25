package main

import "github.com/valyala/fasthttp"

var defaultPool = newWorkerPool(64)

type worker struct {
	count int
	pool  *workerPool
	rq    []*fasthttp.Request
	rs    []*fasthttp.Response
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	p.workers <- w
	return w
}

func (w *worker) reset() {
	w.count = 0
	w.rq = w.rq[:0]
	w.rs = w.rs[:0]
}

func (w *worker) run(n int) {
	for i := 0; i < n; i++ {
		w.rq = append(w.rq, fasthttp.AcquireRequest())
		w.rs = append(w.rs, fasthttp.AcquireResponse())
	}

	// Some logic

	for i := 0; i < n; i++ {
		fasthttp.ReleaseRequest(w.rq[i])
		fasthttp.ReleaseResponse(w.rs[i])
	}
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
