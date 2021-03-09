package main

import "sync"

var bytesPool *sync.Pool

func init() {
	bytesPool = &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, 0x4000)
			return &buf
		},
	}
}
