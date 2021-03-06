package main

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randArray(n int) []int {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = rand.Int()
	}
	return data
}

func TestHeapSort(t *testing.T) {
	data := randArray(1024)
	data = heapSort(data)
	if !sort.IntsAreSorted(data) {
		t.Fail()
	}
}

func BenchmarkHeapSortRnd(b *testing.B) {
	n := 1024
	orig := randArray(n)
	data := make([]int, n)
	for i := 0; i < b.N; i++ {
		copy(data, orig)
		heapSort(data)
	}
}

func BenchmarkStandartSortRnd(b *testing.B) {
	n := 1024
	orig := randArray(n)
	data := make([]int, n)
	for i := 0; i < b.N; i++ {
		copy(data, orig)
		sort.Ints(data)
	}
}
