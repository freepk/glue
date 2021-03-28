package main

type worker struct {
	pool *workerPool
	join chan []byte
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.join = make(chan []byte)
	p.workers <- w
	return w
}

func makeHttpRequest(buf []byte, url string) []byte {
	return buf
}

func makeHttpRequestAsync(buf []byte, url string, out chan<- []byte) {
	out <- makeHttpRequest(buf, url)
}

func (w *worker) doAsync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		go makeHttpRequestAsync(nil, urls[i], w.join)
	}
	for i := 0; i < n; i++ {
		body := <-w.join
		buf = append(buf, body...)
	}
	return buf
}

func (w *worker) doSync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		body := makeHttpRequest(nil, urls[i])
		buf = append(buf, body...)
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
