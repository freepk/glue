package main

import (
	"github.com/valyala/fasthttp"
	"sync"
)

var client = &fasthttp.Client{
	ReadBufferSize:  1 << 14,
	WriteBufferSize: 1 << 14,
}

type request struct {
	req *fasthttp.Request
	res *fasthttp.Response
}

func (r *request) respStatus() int {
	return r.res.StatusCode()
}

func (r *request) respBody() []byte {
	return r.res.Body()
}

type worker struct {
	pool *workerPool
	join *sync.WaitGroup
	reqs [64]request
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.join = new(sync.WaitGroup)
	p.workers <- w
	return w
}

func doWithDone(r *request, join *sync.WaitGroup, url string) {
	defer join.Done()
	r.req.SetRequestURI(url)
	client.Do(r.req, r.res)
}

func (w *worker) run(buf []byte, urls []string) []byte {
	n := len(urls)
	w.join.Add(n)
	for i := 0; i < n; i++ {
		r := &w.reqs[i]
		r.req = fasthttp.AcquireRequest()
		r.res = fasthttp.AcquireResponse()
		go doWithDone(r, w.join, urls[i])
	}
	w.join.Wait()
	for i := 0; i < n; i++ {
		r := &w.reqs[i]
		buf = append(buf, r.respBody()...)
		fasthttp.ReleaseRequest(r.req)
		fasthttp.ReleaseResponse(r.res)
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
