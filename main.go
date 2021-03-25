package main

import (
	"github.com/valyala/fasthttp"
)

const (
	samplePort  = `:8080`
	servicePort = `:8081`
)

// optimize pool size with connection pool size
var defaultPool = newWorkerPool(128)

func sampleHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString(`{}`)
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	w := defaultPool.acquire()
	defer w.release()
	buf := w.run(nil, []string{
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
		`http://localhost:8080/`,
	})
	ctx.Write(buf)
}

func main() {
	go fasthttp.ListenAndServe(samplePort, sampleHandler)
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
