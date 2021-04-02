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
	remoteShardPrefix = `http://catalog-backend-part`
)

type productResult struct {
	join      *sync.WaitGroup
	service   *productService
	baseQuery []byte
	products  []int
}

func newProductResult(service *productService) *productResult {
	result := new(productResult)
	result.service = service
	result.join = new(sync.WaitGroup)
	return result
}

func (pr *productResult) setBaseQuery(path, args string) {
	pr.baseQuery = pr.baseQuery[:0]
	pr.baseQuery = append(pr.baseQuery, path...)
	pr.baseQuery = append(pr.baseQuery, questMarkChar)
	pr.baseQuery = append(pr.baseQuery, args...)
}

func (pr *productResult) setProductBytes(buf []byte) {
	pr.products, _ = parseUints(pr.products[:0], buf)
	heapSort(pr.products)
	pr.products = dedupInts(pr.products)
}

func (pr *productResult) buildRequestURI(buf []byte, shard int, products []int) []byte {
	buf = pr.service.shardHost(buf, shard)
	buf = append(buf, pr.baseQuery...)
	buf = append(buf, ampersandChar)
	buf = append(buf, remoteProductsArg...)
	buf = append(buf, equalSignChar)
	for i := 0; i < len(products); i++ {
		buf = appendUint(buf, products[i])
		buf = append(buf, semicolonChar)
	}
	return buf
}

func (pr *productResult) requestByShard(shard int, products []int) {

	defer pr.join.Done()

	requestURI := make([]byte, 0, 1024)
	requestURI = pr.buildRequestURI(requestURI, shard, products)

	req := fasthttp.AcquireRequest()
	req.SetRequestURIBytes(requestURI)
	resp := fasthttp.AcquireResponse()

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(err)
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func (pr *productResult) dispatchRequests() {

	i := 0
	j := 0

	for i < len(pr.products) {
		j = i
		shard := productShard(pr.products[i])
		for j < len(pr.products) {
			if shard != productShard(pr.products[j]) {
				break
			}
			j++
		}
		pr.join.Add(1)
		go pr.requestByShard(shard, pr.products[i:j])
		i = j
	}

}

func (pr *productResult) request(path, args string, products []byte) {
	pr.setBaseQuery(path, args)
	pr.setProductBytes(products)
	pr.dispatchRequests()
	pr.join.Wait()
}

type productService struct {
	resultPool *sync.Pool
	baseHost   []byte
}

func newProductService(baseHost string) *productService {
	service := new(productService)
	resultPool := new(sync.Pool)
	resultPool.New = func() interface{} {
		return newProductResult(service)
	}
	service.resultPool = resultPool
	service.setBaseHost(baseHost)
	return service
}

func (svc *productService) setBaseHost(baseHost string) {
	svc.baseHost = append(svc.baseHost[:0], baseHost...)
}

func (svc *productService) shardHost(buf []byte, shard int) []byte {
	buf = append(buf, remoteShardPrefix...)
	buf = appendUint(buf, shard)
	buf = append(buf, dotChar)
	buf = append(buf, svc.baseHost...)
	return buf
}

func productShard(i int) int {
	return i / 1000000
}

func (svc *productService) handleProducts(ctx *fasthttp.RequestCtx, remotePath string) {

	result := svc.resultPool.Get().(*productResult)
	defer svc.resultPool.Put(result)

	args := ctx.QueryArgs()
	products := args.Peek(productsArg)
	args.Del(productsArg)

	result.request(remotePath, args.String(), products)

}

func main() {
	baseHost := `wbx-ru.svc.k8s.dataline`
	addr := `:8080`

	service := newProductService(baseHost)
	handler := func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			switch string(ctx.Path()) {
			case catalogPath, remoteCatalogPath:
				service.handleProducts(ctx, remoteCatalogPath)
			case basketPath, remoteBasketPath:
				service.handleProducts(ctx, remoteBasketPath)
			case filtersPath:
			case metricsPath:
			case statePath:
			}
		}
	}

	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
