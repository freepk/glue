package main

func isEqual(a, b []int) bool {
	size := len(a)
	if size != len(b) {
		return false
	}
	for i := 0; i < size; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

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
	i := 0
	j := 0
	for i < n {
		if (i + 1) < n {
			if r[i] == r[i+1] {
				i++
				continue
			}
		}
		r[j] = r[i]
		i++
		j++
	}
	return r[:j]
}
