package main

func heapSort(a []int) {
	right := len(a)
	if right < 2 {
		return
	}
	left := (right - 2) / 2
	for left >= 0 {
		siftDown(a, left, right)
		left--
	}
	right--
	for right > 0 {
		a[0], a[right] = a[right], a[0]
		siftDown(a, 0, right)
		right--
	}
}

func siftDown(a []int, left, right int) {
	curr := (left * 2) + 1
	next := curr + 1
	for curr < right {
		if next < right {
			if a[curr] < a[next] {
				curr++
				next++
			}
		}
		if a[left] >= a[curr] {
			return
		}
		a[left], a[curr] = a[curr], a[left]
		left = curr
		curr = (left * 2) + 1
		next = curr + 1
	}
}
