package main

func heapSort(a []int) {
	hi := len(a)
	if hi < 2 {
		return
	}
	lo := (hi - 2) / 2
	for lo >= 0 {
		siftDown(a, lo, hi)
		lo--
	}
	hi--
	for hi > 0 {
		a[0], a[hi] = a[hi], a[0]
		siftDown(a, 0, hi)
		hi--
	}
}

func siftDown(a []int, lo, hi int) {
	pos := (lo * 2) + 1
	ext := pos + 1
	for pos < hi {
		if ext < hi {
			if a[pos] < a[ext] {
				pos++
				ext++
			}
		}
		if a[lo] >= a[pos] {
			return
		}
		a[lo], a[pos] = a[pos], a[lo]
		lo = pos
		pos = (lo * 2) + 1
		ext = pos + 1
	}
}
