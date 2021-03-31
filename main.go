package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"unsafe"
)

const (
	catalogPath = `/enrichment/v1/api`
	basketPath  = `/basket/v1/api`
	filtersPath = `/v2/filters/only`
	metricsPath = `/metrics`
	statePath   = `/state`

	productsArg = `nm`
	productsMax = 1024

	remoteCatalogPath = `/catalog`
	remoteBasketPath  = `/basket`
	remoteProductsArg = `product`

	remoteService = `catalog-backend-part`
	remoteDomain  = `wbx-ru.svc.k8s.dataline`
)

func partNum(i int) int {
	return i / 1000000
}

func buildRemoteUrl(buf []byte, remotePath string, products []int) []byte {
	buf = append(buf, `http://`...)
	buf = append(buf, remoteService...)
	buf = appendUint(buf, partNum(products[0]))
	buf = append(buf, dotChar)
	buf = append(buf, remoteDomain...)
	buf = append(buf, remotePath...)
	buf = append(buf, questMarkChar)
	buf = append(buf, `locale=ru&lang=ru&`...)
	buf = append(buf, remoteProductsArg...)
	buf = append(buf, equalSignChar)
	for i := 0; i < len(products); i++ {
		buf = appendUint(buf, products[i])
		buf = append(buf, semicolonChar)
	}
	return buf
}

func asyncCallback(output, input []byte) []byte {
	var err error
	str := *(*string)(unsafe.Pointer(&input))
	if _, output, err = fasthttp.Get(output, str); err != nil {
		log.Print(str, err)
	}
	return output
}

var resultsPool = newLimitedPool(func() interface{} {
	return newAsyncResult(asyncCallback)
}, 128)

func handleProducts(ctx *fasthttp.RequestCtx, remotePath string) {
	args := ctx.QueryArgs()
	products := make([]int, 0, productsMax)
	products, _ = parseUints(products, args.Peek(productsArg))
	heapSort(products)
	products = dedupInts(products)
	result := resultsPool.get().(*asyncResult)
	result.resetTasks()
	defer resultsPool.put(result)
	n := len(products)
	i := 0
	for i < n {
		p := partNum(products[i])
		j := i
		for j < n {
			if partNum(products[j]) != p {
				break
			}
			j++
		}
		task := result.newTask()
		task.appendInput(buildRemoteUrl(task.resetInput(), remotePath, products[i:j]))
		i = j
	}
	body := ctx.Response.Body()
	body = result.await(body)
	ctx.SetBody(body)
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsGet() {
		switch string(ctx.Path()) {
		case catalogPath:
			handleProducts(ctx, remoteCatalogPath)
		case basketPath:
			handleProducts(ctx, remoteBasketPath)
		case filtersPath:
		case metricsPath:
		case statePath:
		}
	}
}

func main() {
	serviceAddr := `:8080`
	fasthttp.ListenAndServe(serviceAddr, serviceHandler)
}
