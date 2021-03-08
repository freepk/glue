package main

import (
	"github.com/valyala/fasthttp"
)

const (
	srcArgName  = `nm`
	dstArgName  = `product`
	servicePort = `:8080`
)

func appendHost(buf []byte, part int) []byte {
	buf = append(buf, `catalog-backend-part`...)
	buf = appendUint(buf, part)
	buf = append(buf, `.wbx-ru.svc.k8s.dataline`...)
	return buf
}

func appendPath(buf, path []byte) []byte {
	return append(buf, path...)
}

func appendArgs(buf []byte, args *fasthttp.Args) []byte {
	args.VisitAll(func(k, v []byte) {
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
	})
	return buf
}

func partNum(n int) int {
	return n / 1000000
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	buf := args.Peek(srcArgName)
	if len(buf) == 0 {
		return
	}
	items := make([]int, 0, 256)
	items, scanLen := parseUints(items, buf)
	if scanLen != len(buf) {
		return
	}
	heapSort(items)
	items = dedupInts(items)
	parts := make([]int, 0, 256)
	for i := 0; i < len(items); i++ {
		parts = append(parts, partNum(items[i]))
	}
	parts = dedupInts(parts)
	args.Del(srcArgName)
	i := 0
	for p := 0; p < len(parts); p++ {
		buf = buf[:0]
		buf = appendHost(buf, parts[p])
		buf = appendPath(buf, ctx.Path())
		buf = appendArgs(buf, args)
		for i < len(items) && partNum(items[i]) == parts[p] {
			buf = appendUint(buf, items[i])
			i++
		}
	}
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
