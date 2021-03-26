package main

import (
	"github.com/valyala/fasthttp"
	"time"
)

const (
	samplePort  = `:8080`
	servicePort = `:8081`
)

var defaultPool = newWorkerPool(256)

var sampleUrls = []string{
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
	`http://localhost:8080/`,
}

func sampleHandler(ctx *fasthttp.RequestCtx) {
	time.Sleep(time.Millisecond * 1)
	ctx.WriteString(`{}`)
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	w := defaultPool.acquire()
	defer w.release()
	body := ctx.Response.Body()
	body = w.doAsync(body, sampleUrls)
	ctx.SetBody(body)
}

func main() {
	go fasthttp.ListenAndServe(samplePort, sampleHandler)
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
