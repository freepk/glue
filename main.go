package main

import (
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

const (
	servicePort = `:8080`
)

var sampleItems = []int{2186603, 2654554, 2968009, 3723312, 4144387, 4588641, 4815118, 4861382, 5024530,
	5146236, 5613752, 5616296, 6170054, 6590837, 6671965, 6790848, 6796776, 7601411, 8139643, 8781920,
	9024198, 9322369, 9434636, 9434637, 9434638, 10004476, 10143840, 10219059, 10459579, 10667312,
	10854001, 10854002, 11088253, 11481196, 11857229, 12016099, 12016103, 12423144, 12507082, 12888079,
	12968083, 13227933, 13595361, 13615120, 13615122, 13807725, 13918691, 13944386, 13978470, 14136935,
	14324471, 14878237, 15025349, 15069435, 15154431, 15163808, 15435674, 15457224, 15556061, 15556062,
	16023989, 16379986, 16780324, 16889371, 17367533, 18029960, 18247707, 18362879, 18565187, 18622848}

func partNum(n int) int {
	return n / 1000000
}

func makeRequest(items []int, args []byte) []byte {
	parts := make([]int, 0, 256)
	parts = parts[:len(items)]
	for i := 0; i < len(items); i++ {
		parts = append(parts, partNum(items[i]))
	}
	parts = dedupInts(parts)
	resources := make([]*[]byte, 0, 256)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(parts))
	i := 0
	for p := 0; p < len(parts); p++ {
		pBuf := bytesPool.Get().(*[]byte)
		resources = append(resources, pBuf)
		for (i < len(items)) && (parts[p] <= items[i]) {
			i++
		}
		workerPool.Run(func() {
			defer waitGroup.Done()
			time.Sleep(100 * time.Millisecond)
		})
	}
	waitGroup.Wait()
	for p := 0; p < len(parts); p++ {
		bytesPool.Put(resources[p])
	}
	return nil
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	makeRequest(sampleItems, []byte(`lang=ru&locale=ru`))
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
