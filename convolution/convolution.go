package convolution

import (
	"math/bits"
	"sync"

	"ac-library-go/modint"
)

func countrZero(n uint32) int {
	return bits.TrailingZeros32(n)
}

func bitCeil(n uint32) uint32 {
	x := uint32(1)
	for x < n {
		x *= 2
	}
	return x
}

func powMod(x, n int64, m int) int64 {
	if m == 1 {
		return 0
	}
	_m := uint32(m)
	r := uint64(1)
	y := uint64(x % int64(m))
	if x < 0 {
		y = uint64((x%int64(m) + int64(m)) % int64(m))
	}
	nn := uint64(n)
	for nn > 0 {
		if nn&1 == 1 {
			r = (r * y) % uint64(_m)
		}
		y = (y * y) % uint64(_m)
		nn >>= 1
	}
	return int64(r)
}

func primitiveRoot(m int) int {
	if m == 2 {
		return 1
	}
	if m == 167772161 {
		return 3
	}
	if m == 469762049 {
		return 3
	}
	if m == 754974721 {
		return 11
	}
	if m == 998244353 {
		return 3
	}

	divs := make([]int, 20)
	divs[0] = 2
	cnt := 1
	x := (m - 1) / 2
	for x%2 == 0 {
		x /= 2
	}
	for i := 3; int64(i)*int64(i) <= int64(x); i += 2 {
		if x%i == 0 {
			divs[cnt] = i
			cnt++
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		divs[cnt] = x
		cnt++
	}

	for g := 2; ; g++ {
		ok := true
		for i := 0; i < cnt; i++ {
			if powMod(int64(g), int64((m-1)/divs[i]), m) == 1 {
				ok = false
				break
			}
		}
		if ok {
			return g
		}
	}
}

type fftInfoStruct[M modint.Modulus] struct {
	root   []modint.StaticModint[M]
	iroot  []modint.StaticModint[M]
	rate2  []modint.StaticModint[M]
	irate2 []modint.StaticModint[M]
	rate3  []modint.StaticModint[M]
	irate3 []modint.StaticModint[M]
}

var fftInfoCache sync.Map

func getFFTInfo[M modint.Modulus]() *fftInfoStruct[M] {
	var m M
	mod := m.Mod()
	if v, ok := fftInfoCache.Load(mod); ok {
		return v.(*fftInfoStruct[M])
	}
	info := computeFFTInfo[M]()
	fftInfoCache.Store(mod, info)
	return info
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func computeFFTInfo[M modint.Modulus]() *fftInfoStruct[M] {
	var m M
	mod := m.Mod()
	g := primitiveRoot(int(mod))
	rank2 := countrZero(mod - 1)

	info := &fftInfoStruct[M]{
		root:   make([]modint.StaticModint[M], rank2+1),
		iroot:  make([]modint.StaticModint[M], rank2+1),
		rate2:  make([]modint.StaticModint[M], max(0, rank2-2+1)),
		irate2: make([]modint.StaticModint[M], max(0, rank2-2+1)),
		rate3:  make([]modint.StaticModint[M], max(0, rank2-3+1)),
		irate3: make([]modint.StaticModint[M], max(0, rank2-3+1)),
	}

	info.root[rank2] = modint.NewStatic[M](int64(g)).Pow(int64((mod - 1) >> rank2))
	info.iroot[rank2] = info.root[rank2].Inv()
	for i := rank2 - 1; i >= 0; i-- {
		info.root[i] = info.root[i+1].Mul(info.root[i+1])
		info.iroot[i] = info.iroot[i+1].Mul(info.iroot[i+1])
	}

	prod := modint.RawStatic[M](1)
	iprod := modint.RawStatic[M](1)
	for i := 0; i <= rank2-2; i++ {
		info.rate2[i] = info.root[i+2].Mul(prod)
		info.irate2[i] = info.iroot[i+2].Mul(iprod)
		prod = prod.Mul(info.iroot[i+2])
		iprod = iprod.Mul(info.root[i+2])
	}

	prod = modint.RawStatic[M](1)
	iprod = modint.RawStatic[M](1)
	for i := 0; i <= rank2-3; i++ {
		info.rate3[i] = info.root[i+3].Mul(prod)
		info.irate3[i] = info.iroot[i+3].Mul(iprod)
		prod = prod.Mul(info.iroot[i+3])
		iprod = iprod.Mul(info.root[i+3])
	}

	return info
}

func butterfly[M modint.Modulus](a []modint.StaticModint[M]) {
	n := len(a)
	h := countrZero(uint32(n))
	info := getFFTInfo[M]()
	var m M
	mod2 := uint64(m.Mod()) * uint64(m.Mod())

	len_ := 0
	for len_ < h {
		if h-len_ == 1 {
			p := 1 << (h - len_ - 1)
			rot := modint.RawStatic[M](1)
			for s := 0; s < (1 << len_); s++ {
				offset := s << (h - len_)
				for i := 0; i < p; i++ {
					l := a[i+offset]
					r := a[i+offset+p].Mul(rot)
					a[i+offset] = l.Add(r)
					a[i+offset+p] = l.Sub(r)
				}
				if s+1 != (1 << len_) {
					rot = rot.Mul(info.rate2[countrZero(^uint32(s))])
				}
			}
			len_++
		} else {
			p := 1 << (h - len_ - 2)
			rot := modint.RawStatic[M](1)
			imag := info.root[2]
			for s := 0; s < (1 << len_); s++ {
				rot2 := rot.Mul(rot)
				rot3 := rot2.Mul(rot)
				offset := s << (h - len_)
				for i := 0; i < p; i++ {
					a0 := uint64(a[i+offset].Val())
					a1 := uint64(a[i+offset+p].Val()) * uint64(rot.Val())
					a2 := uint64(a[i+offset+2*p].Val()) * uint64(rot2.Val())
					a3 := uint64(a[i+offset+3*p].Val()) * uint64(rot3.Val())

					a1na3imag := uint64(modint.RawStatic[M](uint32((a1 + mod2 - a3) % uint64(m.Mod()))).Mul(imag).Val())
					na2 := mod2 - a2

					a[i+offset] = modint.RawStatic[M](uint32((a0 + a2 + a1 + a3) % uint64(m.Mod())))
					a[i+offset+p] = modint.RawStatic[M](uint32((a0 + a2 + (2*mod2 - (a1 + a3))) % uint64(m.Mod())))
					a[i+offset+2*p] = modint.RawStatic[M](uint32((a0 + na2 + a1na3imag) % uint64(m.Mod())))
					a[i+offset+3*p] = modint.RawStatic[M](uint32((a0 + na2 + (mod2 - a1na3imag)) % uint64(m.Mod())))
				}
				if s+1 != (1 << len_) {
					rot = rot.Mul(info.rate3[countrZero(^uint32(s))])
				}
			}
			len_ += 2
		}
	}
}

func butterflyInv[M modint.Modulus](a []modint.StaticModint[M]) {
	n := len(a)
	h := countrZero(uint32(n))
	info := getFFTInfo[M]()
	var m M

	len_ := h
	for len_ > 0 {
		if len_ == 1 {
			p := 1 << (h - len_)
			irot := modint.RawStatic[M](1)
			for s := 0; s < (1 << (len_ - 1)); s++ {
				offset := s << (h - len_ + 1)
				for i := 0; i < p; i++ {
					l := a[i+offset]
					r := a[i+offset+p]
					a[i+offset] = l.Add(r)

					val := uint64(l.Val()) + uint64(m.Mod()) - uint64(r.Val())
					a[i+offset+p] = modint.RawStatic[M](uint32(val % uint64(m.Mod()))).Mul(irot)
				}
				if s+1 != (1 << (len_ - 1)) {
					irot = irot.Mul(info.irate2[countrZero(^uint32(s))])
				}
			}
			len_--
		} else {
			p := 1 << (h - len_)
			irot := modint.RawStatic[M](1)
			iimag := info.iroot[2]
			for s := 0; s < (1 << (len_ - 2)); s++ {
				irot2 := irot.Mul(irot)
				irot3 := irot2.Mul(irot)
				offset := s << (h - len_ + 2)
				for i := 0; i < p; i++ {
					a0 := uint64(a[i+offset].Val())
					a1 := uint64(a[i+offset+p].Val())
					a2 := uint64(a[i+offset+2*p].Val())
					a3 := uint64(a[i+offset+3*p].Val())

					a2na3iimag := uint64(modint.RawStatic[M](uint32((uint64(m.Mod()) + a2 - a3) % uint64(m.Mod()))).Mul(iimag).Val())

					a[i+offset] = modint.RawStatic[M](uint32((a0 + a1 + a2 + a3) % uint64(m.Mod())))

					a[i+offset+p] = modint.RawStatic[M](uint32((a0 + uint64(m.Mod()) - a1 + a2na3iimag) % uint64(m.Mod()))).Mul(irot)

					a[i+offset+2*p] = modint.RawStatic[M](uint32((a0 + a1 + uint64(m.Mod()) - a2 + uint64(m.Mod()) - a3) % uint64(m.Mod()))).Mul(irot2)

					a[i+offset+3*p] = modint.RawStatic[M](uint32((a0 + uint64(m.Mod()) - a1 + uint64(m.Mod()) - a2na3iimag) % uint64(m.Mod()))).Mul(irot3)
				}
				if s+1 != (1 << (len_ - 2)) {
					irot = irot.Mul(info.irate3[countrZero(^uint32(s))])
				}
			}
			len_ -= 2
		}
	}
}

func convolutionNaive[M modint.Modulus](a, b []modint.StaticModint[M]) []modint.StaticModint[M] {
	n := len(a)
	m := len(b)
	ans := make([]modint.StaticModint[M], n+m-1)
	if n < m {
		for j := 0; j < m; j++ {
			for i := 0; i < n; i++ {
				ans[i+j] = ans[i+j].Add(a[i].Mul(b[j]))
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				ans[i+j] = ans[i+j].Add(a[i].Mul(b[j]))
			}
		}
	}
	return ans
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Convolution computes the convolution of two arrays of StaticModint.
func Convolution[M modint.Modulus](a, b []modint.StaticModint[M]) []modint.StaticModint[M] {
	n := len(a)
	m := len(b)
	if n == 0 || m == 0 {
		return []modint.StaticModint[M]{}
	}

	var mod M
	z := bitCeil(uint32(n + m - 1))
	if (mod.Mod()-1)%z != 0 {
		panic("mod-1 is not divisible by the required power of 2")
	}

	if min(n, m) <= 60 {
		return convolutionNaive(a, b)
	}

	a2 := make([]modint.StaticModint[M], z)
	copy(a2, a)
	b2 := make([]modint.StaticModint[M], z)
	copy(b2, b)

	butterfly(a2)
	butterfly(b2)
	for i := uint32(0); i < z; i++ {
		a2[i] = a2[i].Mul(b2[i])
	}
	butterflyInv(a2)

	a2 = a2[:n+m-1]
	iz := modint.NewStatic[M](int64(z)).Inv()
	for i := 0; i < n+m-1; i++ {
		a2[i] = a2[i].Mul(iz)
	}
	return a2
}

func invGcd(a, b int64) int64 {
	a %= b
	if a < 0 {
		a += b
	}
	if a == 0 {
		return 0
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
	return m0
}

type Mod754974721 struct{}

func (Mod754974721) Mod() uint32 { return 754974721 }

type Mod167772161 struct{}

func (Mod167772161) Mod() uint32 { return 167772161 }

type Mod469762049 struct{}

func (Mod469762049) Mod() uint32 { return 469762049 }

// ConvolutionLL computes the convolution of two integer arrays.
func ConvolutionLL(a, b []int64) []int64 {
	n := len(a)
	m := len(b)
	if n == 0 || m == 0 {
		return []int64{}
	}

	const MOD1_C = 754974721
	const MOD2_C = 167772161
	const MOD3_C = 469762049
	MOD1 := uint64(MOD1_C)
	MOD2 := uint64(MOD2_C)
	MOD3 := uint64(MOD3_C)
	M2M3 := MOD2 * MOD3
	M1M3 := MOD1 * MOD3
	M1M2 := MOD1 * MOD2
	M1M2M3 := M1M2 * MOD3

	i1 := invGcd(int64(MOD2_C)*int64(MOD3_C), int64(MOD1_C))
	i2 := invGcd(int64(MOD1_C)*int64(MOD3_C), int64(MOD2_C))
	i3 := invGcd(int64(MOD1_C)*int64(MOD2_C), int64(MOD3_C))

	a1 := make([]modint.StaticModint[Mod754974721], n)
	b1 := make([]modint.StaticModint[Mod754974721], m)
	for i := 0; i < n; i++ {
		a1[i] = modint.NewStatic[Mod754974721](a[i])
	}
	for i := 0; i < m; i++ {
		b1[i] = modint.NewStatic[Mod754974721](b[i])
	}
	c1 := Convolution(a1, b1)

	a2 := make([]modint.StaticModint[Mod167772161], n)
	b2 := make([]modint.StaticModint[Mod167772161], m)
	for i := 0; i < n; i++ {
		a2[i] = modint.NewStatic[Mod167772161](a[i])
	}
	for i := 0; i < m; i++ {
		b2[i] = modint.NewStatic[Mod167772161](b[i])
	}
	c2 := Convolution(a2, b2)

	a3 := make([]modint.StaticModint[Mod469762049], n)
	b3 := make([]modint.StaticModint[Mod469762049], m)
	for i := 0; i < n; i++ {
		a3[i] = modint.NewStatic[Mod469762049](a[i])
	}
	for i := 0; i < m; i++ {
		b3[i] = modint.NewStatic[Mod469762049](b[i])
	}
	c3 := Convolution(a3, b3)

	ans := make([]int64, n+m-1)
	offset := []uint64{0, 0, M1M2M3, 2 * M1M2M3, 3 * M1M2M3}
	for i := 0; i < n+m-1; i++ {
		x := uint64(0)
		x += (uint64(c1[i].Val()) * uint64(i1)) % MOD1 * M2M3
		x += (uint64(c2[i].Val()) * uint64(i2)) % MOD2 * M1M3
		x += (uint64(c3[i].Val()) * uint64(i3)) % MOD3 * M1M2
		diff := int64(c1[i].Val()) - int64(x%MOD1)
		if diff < 0 {
			diff += int64(MOD1_C)
		}
		x -= offset[diff%5]
		ans[i] = int64(x)
	}
	return ans
}
