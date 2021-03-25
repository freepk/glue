package main

import (
	"github.com/valyala/fasthttp"
	"sync"
)

type request struct {
	req *fasthttp.Request
	res *fasthttp.Response
}

func newRequest() *request {
	r := new(request)
	r.req = fasthttp.AcquireRequest()
	r.res = fasthttp.AcquireResponse()
	return r
}

func (r *request) doWithDone(join *sync.WaitGroup, url string) {
	defer join.Done()
	r.req.SetRequestURI(url)
	fasthttp.Do(r.req, r.res)
}

func (r *request) respStatus() int {
	return r.res.StatusCode()
}

func (r *request) respBody() []byte {
	return r.res.Body()
}

func (r *request) release() {
	fasthttp.ReleaseRequest(r.req)
	fasthttp.ReleaseResponse(r.res)
}

type worker struct {
	pool *workerPool
	join *sync.WaitGroup
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.join = new(sync.WaitGroup)
	p.workers <- w
	return w
}

func (w *worker) run(buf []byte, urls []string) []byte {
	reqs := make([]*request, 0, 32)
	n := len(urls)
	w.join.Add(n)
	for i := 0; i < n; i++ {
		r := newRequest()
		reqs = append(reqs, r)
		go r.doWithDone(w.join, urls[i])
	}
	w.join.Wait()
	for i := 0; i < n; i++ {
		r := reqs[i]
		buf = append(buf, r.respBody()...)
		r.release()
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
	return w
}
