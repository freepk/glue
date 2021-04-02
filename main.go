package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"sync"
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

type productResult struct {
	wait     *sync.WaitGroup
	products []int
}

func newProductResult() *productResult {
	return &productResult{
		wait: &sync.WaitGroup{},
	}
}

func (pr *productResult) setProductBytes(buf []byte) {
	pr.products, _ = parseUints(pr.products[:0], buf)
	heapSort(pr.products)
	pr.products = dedupInts(pr.products)
}

func (pr *productResult) requestShardProducts(shard int, products []int) {
	defer pr.wait.Done()
	// fmt.Println(shard, products)
}

func (pr *productResult) requestProducts() {
	n := len(pr.products)
	i := 0
	j := 0
	for i < n {
		j = i
		shard := productShard(pr.products[i])
		for j < n {
			if shard != productShard(pr.products[j]) {
				break
			}
			j++
		}
		pr.wait.Add(1)
		go pr.requestShardProducts(shard, pr.products[i:j])
		i = j
	}
	pr.wait.Wait()
}

type productService struct {
	shardURLs  []string
	resultPool *sync.Pool
}

func newService() *productService {
	return &productService{
		resultPool: &sync.Pool{
			New: func() interface{} {
				return newProductResult()
			},
		},
	}
}

func productShard(i int) int {
	return i / 1000000
}

func (svc *productService) handleProducts(ctx *fasthttp.RequestCtx, remotePath string) {

	result := svc.resultPool.Get().(*productResult)
	defer svc.resultPool.Put(result)

	result.setProductBytes(ctx.QueryArgs().Peek(productsArg))
	result.requestProducts()

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
