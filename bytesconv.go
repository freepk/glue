package main

const (
	ampersandChar = 0x26
	dotChar       = 0x2e
	zeroChar      = 0x30
	semicolonChar = 0x3b
	equalSignChar = 0x3d
	questMarkChar = 0x3f
)

func parseUint(b []byte) (int, int) {
	n := len(b)
	i := 0
	r := 0
	c := byte(0)
	for i < n {
		c = b[i] - zeroChar
		if c > 9 {
			break
		}
		r *= 10
		r += int(c)
		i++
	}
	return r, i
}

func appendUint(r []byte, n int) []byte {
	b := make([]byte, 32)
	i := len(b)
	m := 0
	for n > 9 {
		i--
		m = n
		n = n / 10
		m -= n * 10
		b[i] = byte(m + zeroChar)
	}
	i--
	b[i] = byte(n + zeroChar)
	r = append(r, b[i:]...)
	return r
}

func parseUints(r []int, b []byte) ([]int, int) {
	n := len(b)
	i := 0
	for i < n {
		if b[i] == semicolonChar {
			i++
			continue
		}
		x, j := parseUint(b[i:])
		if j == 0 {
			break
		}
		r = append(r, x)
		i += j
	}
	return r, i
}

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

	i := 0
	j := 0

	for i < len(a) {
		j = i
		p := fn(a[i])
		for j < len(a) {
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
