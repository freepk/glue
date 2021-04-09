package main

func dedupInts(r []int) []int {
	n := len(r)
	if n < 2 {
		return r
	}
	i := 0
	j := 1
	for j < n {
		if r[i] != r[j] {
			i++
			r[i] = r[j]
		}
		j++
	}
	i++
	return r[:i]
}

func splitByFunc(r [][]int, a []int, fn func(int) int) [][]int {
	n := len(a)
	i := 0
	j := 0
	for i < n {
		j = i
		p := fn(a[i])
		for j < n {
			if p != fn(a[j]) {
				break
			}
			j++
		}
		r = append(r, a[i:j])
		i = j
	}
	return r
}
