package main

func scanNum(b []byte) (int, int) {
	n := len(b)
	i := 0
	r := 0
	for i < n {
		if b[i] < 0x30 || b[i] > 0x39 {
			break
		}
		r *= 10
		r += int(b[i]) - 0x30
		i++
	}
	return r, i
}

func scanNums(r []int, b []byte) ([]int, int) {
	n := len(b)
	i := 0
	for i < n {
		if b[i] == 0x3b {
			i++
			continue
		}
		x, j := scanNum(b[i:])
		if j == 0 {
			break
		}
		r = append(r, x)
		i += j
	}
	return r, i
}

func dedupNums(r []int) []int {
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
