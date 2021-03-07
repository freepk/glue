package main

import (
	//"fmt"
	"github.com/valyala/fasthttp"
)

const (
	argName     = `nm`
	servicePort = `:8080`
)

func serviceHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	buf := args.Peek(argName)
	bufLen := len(buf)
	if bufLen == 0 {
		return
	}
	items := make([]int, 0, 256)
	items, scanLen := scanNums(items, buf)
	if scanLen != bufLen {
		return
	}
	//fmt.Println("args", string(buf))
	heapSort(items)
	items = dedupNums(items)
	//fmt.Println("items", items)
	shards := make([]int, len(items))
	for i, item := range items {
		shards[i] = item / 1000000
	}
	shards = dedupNums(shards)
	//fmt.Println("shards", shards)
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
