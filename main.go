package main

import (
	"github.com/valyala/fasthttp"
)

const (
	argumentName = `nm`
	servicePort  = `:8080`
)

func serviceHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	buff := args.Peek(argumentName)
	buffLen := len(buff)
	if buffLen == 0 {
		return
	}
	items := make([]int, 0, 256)
	items, scanLen := scanNums(items, buff)
	if scanLen != buffLen {
		return
	}
	heapSort(items)
	shards := make([]int, 0, 256)
	for i := 0; i < len(items); i++ {
		shards = append(shards, (items[i] / 1000000))
	}
	shards = dedupNums(shards)
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
