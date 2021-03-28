package main

type worker struct {
	pool *workerPool
	join chan []byte
	data [][]byte
}

func newWorker(p *workerPool) *worker {
	w := new(worker)
	w.pool = p
	w.join = make(chan []byte)
	w.data = make([][]byte, 256)
	p.workers <- w
	return w
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
		w.data[i] = makeHttpRequest(w.data[i][:0], urls[i])
		buf = append(buf, w.data[i]...)
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
