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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Int()
	}
	return a
}

func TestHeapSort(t *testing.T) {
	a := randArray(1024 * 1024)
	heapSort(a)
	if !sort.IntsAreSorted(a) {
		t.Fail()
	}
}

func BenchmarkHeapSortRnd(b *testing.B) {
	n := 1024
	a := randArray(n)
	c := make([]int, n)
	for i := 0; i < b.N; i++ {
		copy(c, a)
		heapSort(c)
	}
}

func BenchmarkStandartSortRnd(b *testing.B) {
	n := 1024
	a := randArray(n)
	c := make([]int, n)
	for i := 0; i < b.N; i++ {
		copy(c, a)
		sort.Ints(c)
	}
}
