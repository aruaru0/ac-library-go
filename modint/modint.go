package modint

// Modulus is an interface that returns a modulo value.
type Modulus interface {
	Mod() uint32
}

// Mod998244353 represents the modulo 998244353.
type Mod998244353 struct{}

func (Mod998244353) Mod() uint32 { return 998244353 }

// Mod1000000007 represents the modulo 1000000007.
type Mod1000000007 struct{}

func (Mod1000000007) Mod() uint32 { return 1000000007 }

type Mod18433 struct{}
func (Mod18433) Mod() uint32 { return 18433 }

type Mod2 struct{}
func (Mod2) Mod() uint32 { return 2 }

type Mod257 struct{}
func (Mod257) Mod() uint32 { return 257 }

// StaticModint is a modulo integer with a statically defined modulus.
type StaticModint[M Modulus] struct {
	v uint32
}

// NewStatic creates a new StaticModint.
func NewStatic[M Modulus](v int64) StaticModint[M] {
	var m M
	mod := int64(m.Mod())
	v %= mod
	if v < 0 {
		v += mod
	}
	return StaticModint[M]{v: uint32(v)}
}

// RawStatic creates a new StaticModint without taking the modulo.
// It is the caller's responsibility to ensure that v < Mod().
func RawStatic[M Modulus](v uint32) StaticModint[M] {
	return StaticModint[M]{v: v}
}

func (m StaticModint[M]) Val() uint32 {
	return m.v
}

func (m StaticModint[M]) Add(rhs StaticModint[M]) StaticModint[M] {
	var mod M
	v := m.v + rhs.v
	if v >= mod.Mod() {
		v -= mod.Mod()
	}
	return StaticModint[M]{v: v}
}

func (m StaticModint[M]) Sub(rhs StaticModint[M]) StaticModint[M] {
	var mod M
	v := m.v - rhs.v
	if v >= mod.Mod() {
		v += mod.Mod()
	}
	return StaticModint[M]{v: v}
}

func (m StaticModint[M]) Mul(rhs StaticModint[M]) StaticModint[M] {
	var mod M
	v := uint64(m.v) * uint64(rhs.v)
	return StaticModint[M]{v: uint32(v % uint64(mod.Mod()))}
}

func (m StaticModint[M]) Div(rhs StaticModint[M]) StaticModint[M] {
	return m.Mul(rhs.Inv())
}

func (m StaticModint[M]) Pow(n int64) StaticModint[M] {
	if n < 0 {
		panic("n must be non-negative")
	}
	x := m
	r := RawStatic[M](1)
	for n > 0 {
		if n&1 == 1 {
			r = r.Mul(x)
		}
		x = x.Mul(x)
		n >>= 1
	}
	return r
}

func invGCD(a, b int64) (int64, int64) {
	a %= b
	if a < 0 {
		a += b
	}
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

func (m StaticModint[M]) Inv() StaticModint[M] {
	var mod M
	g, x := invGCD(int64(m.v), int64(mod.Mod()))
	if g != 1 {
		panic("gcd(m.v, mod) must be 1")
	}
	return StaticModint[M]{v: uint32(x)}
}

// Aliases for common moduli
type Modint998244353 = StaticModint[Mod998244353]
type Modint1000000007 = StaticModint[Mod1000000007]

func NewModint998244353(v int64) Modint998244353 {
	return NewStatic[Mod998244353](v)
}

func NewModint1000000007(v int64) Modint1000000007 {
	return NewStatic[Mod1000000007](v)
}

// DynamicModint is a modulo integer with a dynamic modulus.
type DynamicModint struct {
	v uint32
	m uint32
}

// NewDynamic creates a new DynamicModint with a specific modulus.
func NewDynamic(v int64, mod uint32) DynamicModint {
	if mod < 1 {
		panic("mod must be at least 1")
	}
	m := int64(mod)
	v %= m
	if v < 0 {
		v += m
	}
	return DynamicModint{v: uint32(v), m: mod}
}

// RawDynamic creates a new DynamicModint without taking the modulo.
func RawDynamic(v uint32, mod uint32) DynamicModint {
	return DynamicModint{v: v, m: mod}
}

func (m DynamicModint) Val() uint32 {
	return m.v
}

func (m DynamicModint) Mod() uint32 {
	return m.m
}

func (m DynamicModint) Add(rhs DynamicModint) DynamicModint {
	if m.m != rhs.m {
		panic("mod mismatch")
	}
	v := m.v + rhs.v
	if v >= m.m {
		v -= m.m
	}
	return DynamicModint{v: v, m: m.m}
}

func (m DynamicModint) Sub(rhs DynamicModint) DynamicModint {
	if m.m != rhs.m {
		panic("mod mismatch")
	}
	v := m.v - rhs.v
	if v >= m.m {
		v += m.m
	}
	return DynamicModint{v: v, m: m.m}
}

func (m DynamicModint) Mul(rhs DynamicModint) DynamicModint {
	if m.m != rhs.m {
		panic("mod mismatch")
	}
	v := uint64(m.v) * uint64(rhs.v)
	return DynamicModint{v: uint32(v % uint64(m.m)), m: m.m}
}

func (m DynamicModint) Div(rhs DynamicModint) DynamicModint {
	return m.Mul(rhs.Inv())
}

func (m DynamicModint) Pow(n int64) DynamicModint {
	if n < 0 {
		panic("n must be non-negative")
	}
	x := m
	r := RawDynamic(1, m.m)
	for n > 0 {
		if n&1 == 1 {
			r = r.Mul(x)
		}
		x = x.Mul(x)
		n >>= 1
	}
	return r
}

func (m DynamicModint) Inv() DynamicModint {
	g, x := invGCD(int64(m.v), int64(m.m))
	if g != 1 {
		panic("gcd(m.v, mod) must be 1")
	}
	return DynamicModint{v: uint32(x), m: m.m}
}
