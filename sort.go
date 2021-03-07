package main

func heapSort(a []int) {
	n := len(a)
	i := (n - 1) / 2
	for i >= 0 {
		siftDown(a, i, n)
		i--
	}
	i = (n - 1)
	for i > 0 {
		a[0], a[i] = a[i], a[0]
		siftDown(a, 0, i)
		i--
	}
}

func siftDown(a []int, i, n int) {
	j := (i * 2) + 1
	k := j + 1
	for j < n {
		if k < n {
			if a[j] < a[k] {
				j++
				k++
			}
		}
		if a[i] >= a[j] {
			break
		}
		a[i], a[j] = a[j], a[i]
		i = j
		j = (i * 2) + 1
		k = j + 1
	}
}
