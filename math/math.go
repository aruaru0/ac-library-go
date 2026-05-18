package math

// safeMod calculates x mod m and always returns a non-negative result.
func safeMod(x, m int64) int64 {
	x %= m
	if x < 0 {
		x += m
	}
	return x
}

// invGCD returns (g, x) such that g = gcd(a, b), x*a = g (mod b), 0 <= x < b/g
func invGCD(a, b int64) (int64, int64) {
	a = safeMod(a, b)
	if a == 0 {
		return b, 0
	}
	s, t := b, a
	m0, m1 := int64(0), int64(1)
	for t != 0 {
		u := s / t
		s -= t * u
		m0 -= m1 * u

		s, t = t, s
		m0, m1 = m1, m0
	}
	if m0 < 0 {
		m0 += b / s
	}
	return s, m0
}

// PowMod returns (x ** n) % m.
func PowMod(x, n int64, m int) int64 {
	if n < 0 {
		panic("n must be non-negative")
	}
	if m < 1 {
		panic("m must be at least 1")
	}
	if m == 1 {
		return 0
	}
	_m := uint64(m)
	r := uint64(1)
	y := uint64(safeMod(x, int64(m)))
	for n > 0 {
		if n&1 == 1 {
			r = (r * y) % _m
		}
		y = (y * y) % _m
		n >>= 1
	}
	return int64(r)
}

// InvMod returns y such that (x * y) % m == 1 (0 <= y < m).
func InvMod(x, m int64) int64 {
	if m < 1 {
		panic("m must be at least 1")
	}
	g, z := invGCD(x, m)
	if g != 1 {
		panic("gcd(x, m) must be 1")
	}
	return z
}

// CRT solves a system of linear congruences.
// It returns (y, z) such that x = r[i] (mod m[i]) is equivalent to x = y (mod z).
// If there is no solution, it returns (0, 0).
func CRT(r, m []int64) (int64, int64) {
	if len(r) != len(m) {
		panic("len(r) must be equal to len(m)")
	}
	n := len(r)
	r0 := int64(0)
	m0 := int64(1)
	for i := 0; i < n; i++ {
		if m[i] < 1 {
			panic("m[i] must be at least 1")
		}
		r1 := safeMod(r[i], m[i])
		m1 := m[i]
		if m0 < m1 {
			r0, r1 = r1, r0
			m0, m1 = m1, m0
		}
		if m0%m1 == 0 {
			if r0%m1 != r1 {
				return 0, 0
			}
			continue
		}
		g, im := invGCD(m0, m1)
		u1 := m1 / g
		if (r1-r0)%g != 0 {
			return 0, 0
		}
		x := (r1 - r0) / g % u1 * im % u1
		r0 += x * m0
		m0 *= u1 // m0 becomes lcm(m0, m1)
		if r0 < 0 {
			r0 += m0
		}
	}
	return r0, m0
}

// floorSumUnsigned is the internal unsigned implementation for FloorSum
func floorSumUnsigned(n, m, a, b uint64) uint64 {
	ans := uint64(0)
	for {
		if a >= m {
			ans += n * (n - 1) / 2 * (a / m)
			a %= m
		}
		if b >= m {
			ans += n * (b / m)
			b %= m
		}
		yMax := a*n + b
		if yMax < m {
			break
		}
		n = yMax / m
		b = yMax % m
		m, a = a, m
	}
	return ans
}

// FloorSum returns sum_{i=0}^{n-1} floor((a * i + b) / m).
func FloorSum(n, m, a, b int64) int64 {
	if n < 0 || n >= (1<<32) {
		panic("n must be in [0, 2^32)")
	}
	if m < 1 || m >= (1<<32) {
		panic("m must be in [1, 2^32)")
	}
	ans := uint64(0)
	if a < 0 {
		a2 := uint64(safeMod(a, m))
		ans -= uint64(n) * uint64(n-1) / 2 * ((a2 - uint64(a)) / uint64(m))
		a = int64(a2)
	}
	if b < 0 {
		b2 := uint64(safeMod(b, m))
		ans -= uint64(n) * ((b2 - uint64(b)) / uint64(m))
		b = int64(b2)
	}
	ans += floorSumUnsigned(uint64(n), uint64(m), uint64(a), uint64(b))
	return int64(ans)
}
