package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

const (
	inArgName   = `nm`
	outArgName  = `product`
	servicePort = `:8080`
)

func partHost(n int) []byte {
	buf := make([]byte, 0, 256)
	buf = append(buf, `catalog-backend-part`...)
	buf = appendUint(buf, n)
	buf = append(buf, `.wbx-ru.svc.k8s.dataline`...)
	return buf
}

func partNum(n int) int {
	return n / 1000000
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	buf := args.Peek(inArgName)
	args.Del(inArgName)
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
	buf = buf[:0]
	args.VisitAll(func(k, v []byte) {
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
	})
	buf = append(buf, outArgName...)
	buf = append(buf, '=')
	i := 0
	for p := 0; p < len(parts); p++ {
		tmp := buf
		for i < len(items) && partNum(items[i]) == parts[p] {
			tmp = appendUint(tmp, items[i])
			tmp = append(tmp, 0x3b)
			i++
		}
		log.Println(string(tmp))
	}
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
