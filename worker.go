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

func doSome() []byte {
	return nil
}

func doSomeAsync(buf []byte, out chan<- []byte) {
	out <- doSome()
}

func (w *worker) doAsync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		go doSomeAsync(nil, w.join)
	}
	for i := 0; i < n; i++ {
		<-w.join
	}
	return buf
}

func (w *worker) doSync(buf []byte, urls []string) []byte {
	n := len(urls)
	for i := 0; i < n; i++ {
		doSome()
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
