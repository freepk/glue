package main

import "github.com/valyala/fasthttp"

func makeHttpRequest(buf []byte, url string) []byte {
	_, buf, _ = fasthttp.Get(buf, url)
	return buf
}

func makeHttpRequestAsync(buf []byte, url string, join chan<- []byte) {
	join <- makeHttpRequest(buf, url)
}
