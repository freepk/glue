package main

func heapSort(a []int) []int {
	heapify(a)
	for i := len(a) - 1; i > 0; i-- {
		a[0], a[i] = a[i], a[0]
		siftDown(a, 0, i)
	}
	return a
}

func heapify(a []int) {
	n := len(a)
	for i := (n - 1) / 2; i >= 0; i-- {
		siftDown(a, i, n)
	}
}

func siftDown(h []int, lo, hi int) {
	r := lo
	for {
		c := (r * 2) + 1
		if c >= hi {
			break
		}
		if ((c + 1) < hi) && (h[c] < h[c+1]) {
			c++
		}
		if h[r] < h[c] {
			h[r], h[c] = h[c], h[r]
			r = c
		} else {
			break
		}
	}
}
