package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

const (
	catalogPath = `/enrichment/v1/api`
	basketPath  = `/basket/v1/api`
	filtersPath = `/v2/filters/only`
	metricsPath = `/metrics`
	statePath   = `/state`
	productsArg = `nm`
)

const (
	remoteCatalogPath = `/catalog`
	remoteBasketPath  = `/basket`
	remoteProductsArg = `product`
)

type service struct {
	shards []string
}

func newService() *service {
	return new(service)
}

func (svc *service) handleProducts(ctx *fasthttp.RequestCtx, remotePath string) {
	products, _ := parseUints(make([]int, 0, 1024), ctx.QueryArgs().Peek(productsArg))
	heapSort(products)
	products = dedupInts(products)

	args := ctx.QueryArgs()
	args.Del(productsArg)

	i := 0
	j := 0
	for i < len(products) {
		j = i
		for j < len(products) {
			if (products[i] / 1000000) != (products[j] / 1000000) {
				break
			}
			j++
		}
		i = j
	}
}

func main() {
	svc := newService()
	handler := func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			switch string(ctx.Path()) {
			case catalogPath, remoteCatalogPath:
				svc.handleProducts(ctx, remoteCatalogPath)
			case basketPath, remoteBasketPath:
				svc.handleProducts(ctx, remoteBasketPath)
			case filtersPath:
			case metricsPath:
			case statePath:
			}
		}
	}
	addr := `:8080`
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
