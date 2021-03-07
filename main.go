package main

import (
	"github.com/valyala/fasthttp"
)

const (
	inputName   = `product`
	servicePort = `:8080`
)

func serviceHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	buf := args.Peek(inputName)
	bufLen := len(buf)
	if bufLen == 0 {
		return
	}
	nums := make([]int, 0, 256)
	// string.Split makes too many string allocations in heap
	nums, scanLen := scanNums(nums, buf)
	if scanLen != bufLen {
		return
	}
	// sort.Ints provocate compilter escape nums to heap
	heapSort(nums)
	nums = dedupNums(nums)
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
