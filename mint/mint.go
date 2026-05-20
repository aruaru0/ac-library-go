package mint

// ModInt represents a modular arithmetic context.
type ModInt struct {
	mod       int
	factMemo  []int
	ifactMemo []int
}

// NewModInt creates a new ModInt context with the given modulus.
func NewModInt(mod int) *ModInt {
	return &ModInt{
		mod:       mod,
		factMemo:  []int{1, 1},
		ifactMemo: []int{1, 1},
	}
}

// Add returns (a + b) % mod.
func (m *ModInt) Add(a, b int) int {
	ret := (a + b) % m.mod
	if ret < 0 {
		ret += m.mod
	}
	return ret
}

// Sub returns (a - b) % mod.
func (m *ModInt) Sub(a, b int) int {
	ret := (a - b) % m.mod
	if ret < 0 {
		ret += m.mod
	}
	return ret
}

// Mul returns (a * b) % mod.
func (m *ModInt) Mul(a, b int) int {
	return int(int64(a) * int64(b) % int64(m.mod))
}

// Div returns (a / b) % mod.
func (m *ModInt) Div(a, b int) int {
	return m.Mul(a, m.ModInv(b))
}

// Pow returns (p ^ n) % mod.
func (m *ModInt) Pow(p, n int) int {
	ret := 1
	x := p % m.mod
	for n != 0 {
		if n%2 == 1 {
			ret = m.Mul(ret, x)
		}
		n /= 2
		x = m.Mul(x, x)
	}
	return ret
}

// ModInv returns the modular inverse of a modulo mod.
func (m *ModInt) ModInv(a int) int {
	b, u, v := m.mod, 1, 0
	for b != 0 {
		t := a / b
		a -= t * b
		a, b = b, a
		u -= t * v
		u, v = v, u
	}
	u %= m.mod
	if u < 0 {
		u += m.mod
	}
	return u
}

// InitComb precomputes factorials and inverse factorials up to limit in O(limit) time.
func (m *ModInt) InitComb(limit int) {
	if len(m.factMemo) > limit {
		return
	}
	fact := make([]int, limit+1)
	ifact := make([]int, limit+1)

	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % int64(m.mod))
	}

	ifact[limit] = m.ModInv(fact[limit])
	for i := limit - 1; i >= 0; i-- {
		ifact[i] = int(int64(ifact[i+1]) * int64(i+1) % int64(m.mod))
	}

	m.factMemo = fact
	m.ifactMemo = ifact
}

// Fact returns n! % mod.
func (m *ModInt) Fact(n int) int {
	if len(m.factMemo) <= n {
		m.InitComb(n * 2)
	}
	return m.factMemo[n]
}

// IFact returns (1 / n!) % mod.
func (m *ModInt) IFact(n int) int {
	if len(m.ifactMemo) <= n {
		m.InitComb(n * 2)
	}
	return m.ifactMemo[n]
}

// NCr returns nCr % mod.
func (m *ModInt) NCr(n, r int) int {
	if n < r || r < 0 {
		return 0
	}
	if len(m.factMemo) <= n {
		m.InitComb(n * 2)
	}
	return m.Mul(m.Fact(n), m.Mul(m.IFact(r), m.IFact(n-r)))
}

// PowModMatrix returns A^p % mod.
func (m *ModInt) PowModMatrix(A [][]int, p int) [][]int {
	N := len(A)
	ret := make([][]int, N)
	for i := 0; i < N; i++ {
		ret[i] = make([]int, N)
		ret[i][i] = 1
	}

	for p > 0 {
		if p&1 == 1 {
			ret = m.MulMod(ret, A)
		}
		A = m.MulMod(A, A)
		p >>= 1
	}
	return ret
}

// MulMod returns (A * B) % mod.
func (m *ModInt) MulMod(A, B [][]int) [][]int {
	H := len(A)
	W := len(B[0])
	K := len(A[0])

	C := make([][]int, H)
	for i := 0; i < H; i++ {
		C[i] = make([]int, W)
	}

	for i := 0; i < H; i++ {
		for j := 0; j < W; j++ {
			var sum int64
			for k := 0; k < K; k++ {
				sum += int64(A[i][k]) * int64(B[k][j])
			}
			C[i][j] = int(sum % int64(m.mod))
		}
	}
	return C
}
