package main

import (
	"github.com/valyala/fasthttp"
	"time"
)

const (
	samplePort = `:8080`
	asyncPort  = `:8081`
	syncPort   = `:8082`
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

func asyncHandler(ctx *fasthttp.RequestCtx) {
	w := defaultPool.acquire()
	defer w.release()
	body := ctx.Response.Body()
	body = w.doAsync(body, sampleUrls)
	ctx.SetBody(body)
}

func syncHandler(ctx *fasthttp.RequestCtx) {
	w := defaultPool.acquire()
	defer w.release()
	body := ctx.Response.Body()
	body = w.doSync(body, sampleUrls)
	ctx.SetBody(body)
}

func main() {
	go fasthttp.ListenAndServe(asyncPort, asyncHandler)
	go fasthttp.ListenAndServe(syncPort, syncHandler)
	fasthttp.ListenAndServe(samplePort, sampleHandler)
}
