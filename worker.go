package main

import (
	"github.com/valyala/fasthttp"
)

type worker struct {
	pool *workerPool
	join chan []byte
	data [][]byte
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.join = make(chan []byte)
	w.data = make([][]byte, 64)
	p.workers <- w
	return w
}

func makeHttpRequest(buf []byte, url string) []byte {
	_, buf, _ = fasthttp.Get(buf, url)
	return buf
}

func makeHttpRequestAsync(buf []byte, url string, join chan<- []byte) {
	join <- makeHttpRequest(buf, url)
}

func (w *worker) doAsync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		go makeHttpRequestAsync(w.data[i][:0], urls[i], w.join)
	}
	for i := 0; i < n; i++ {
		w.data[i] = <-w.join
		buf = append(buf, w.data[i]...)
	}
	return buf
}

func (w *worker) doSync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		buf = makeHttpRequest(buf, urls[i])
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
