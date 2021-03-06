// +build pprof

package main

import "net/http"
import _ "net/http/pprof"

func init() {
	// go tool pprof http://localhost:6060/debug/pprof/profile
	go http.ListenAndServe(":6060", nil)
}
