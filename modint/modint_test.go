package modint

import (
	"math"
	"testing"
)

func TestModintDynamicBorder(t *testing.T) {
	modUpper := uint32(math.MaxInt32)
	for mod := modUpper; mod >= modUpper-20; mod-- {
		v := []int64{}
		for i := int64(0); i < 10; i++ {
			v = append(v, i)
			v = append(v, int64(mod)-i)
			v = append(v, int64(mod)/2+i)
			v = append(v, int64(mod)/2-i)
		}
		for _, a := range v {
			ma := NewDynamic(a, mod)
			a2 := a % int64(mod)
			if a2 < 0 {
				a2 += int64(mod)
			}
			cube := ((a2 * a2) % int64(mod) * a2) % int64(mod)
			if ma.Pow(3).Val() != uint32(cube) {
				t.Errorf("Pow failed for %v", a)
			}
			for _, b := range v {
				mb := NewDynamic(b, mod)
				b2 := b % int64(mod)
				if b2 < 0 {
					b2 += int64(mod)
				}
				add := (a2 + b2) % int64(mod)
				sub := (a2 - b2 + int64(mod)) % int64(mod)
				mul := (a2 * b2) % int64(mod)
				if ma.Add(mb).Val() != uint32(add) {
					t.Errorf("Add failed")
				}
				if ma.Sub(mb).Val() != uint32(sub) {
					t.Errorf("Sub failed")
				}
				if ma.Mul(mb).Val() != uint32(mul) {
					t.Errorf("Mul failed")
				}
			}
		}
	}
}

func TestModintMod1(t *testing.T) {
	for i := int64(0); i < 100; i++ {
		for j := int64(0); j < 100; j++ {
			if NewDynamic(i, 1).Mul(NewDynamic(j, 1)).Val() != 0 {
				t.Errorf("Mod1 Mul failed")
			}
		}
	}
	if NewDynamic(1234, 1).Add(NewDynamic(5678, 1)).Val() != 0 {
		t.Errorf("Mod1 Add failed")
	}
	if NewDynamic(0, 1).Inv().Val() != 0 {
		t.Errorf("Mod1 Inv failed")
	}
}

func TestModintStaticUsage(t *testing.T) {
	a := NewModint1000000007(1)
	b := NewModint1000000007(3)
	if a.Val() == b.Val() {
		t.Errorf("Static equality failed")
	}
	
	c := NewModint998244353(998244353 - 1)
	if c.Add(NewModint998244353(2)).Val() != 1 {
		t.Errorf("Static wrap-around failed")
	}

	d := NewModint998244353(2)
	if d.Pow(10).Val() != 1024 {
		t.Errorf("Static Pow failed")
	}
	
	e := NewModint998244353(2)
	invE := e.Inv()
	if e.Mul(invE).Val() != 1 {
		t.Errorf("Static Inv failed")
	}
}

func TestModintInv(t *testing.T) {
	for i := int64(1); i < 100000; i++ {
		x := NewModint1000000007(i).Inv().Val()
		if (int64(x)*i)%1000000007 != 1 {
			t.Errorf("Static Inv failed for %d", i)
		}
	}

	for i := int64(1); i < 100000; i++ {
		m := NewDynamic(i, 998244353)
		x := m.Inv().Val()
		if (int64(x)*i)%998244353 != 1 {
			t.Errorf("Dynamic Inv failed for %d", i)
		}
	}
}
