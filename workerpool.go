package main

import (
	"github.com/freepk/workerpool"
)

var workerPool *workerpool.Pool

func init() {
	workerPool = workerpool.NewPool(1024)
	go workerPool.Start()
}
