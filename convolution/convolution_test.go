package convolution

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"ac-library-go/modint"
)

func convNaive[M modint.Modulus](a, b []modint.StaticModint[M]) []modint.StaticModint[M] {
	n := len(a)
	m := len(b)
	if n == 0 || m == 0 {
		return []modint.StaticModint[M]{}
	}
	c := make([]modint.StaticModint[M], n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c[i+j] = c[i+j].Add(a[i].Mul(b[j]))
		}
	}
	return c
}

func convLLNaive(a, b []int64) []int64 {
	n := len(a)
	m := len(b)
	if n == 0 || m == 0 {
		return []int64{}
	}
	c := make([]int64, n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c[i+j] += a[i] * b[j]
		}
	}
	return c
}

func TestConvolutionEmpty(t *testing.T) {
	a := []modint.Modint998244353{}
	b := []modint.Modint998244353{}
	if len(Convolution(a, b)) != 0 {
		t.Errorf("Empty failed")
	}

	b2 := []modint.Modint998244353{modint.NewModint998244353(1)}
	if len(Convolution(a, b2)) != 0 {
		t.Errorf("Empty failed")
	}
	if len(Convolution(b2, a)) != 0 {
		t.Errorf("Empty failed")
	}

	all := []int64{}
	bll := []int64{}
	if len(ConvolutionLL(all, bll)) != 0 {
		t.Errorf("Empty LL failed")
	}
	bll2 := []int64{1, 2}
	if len(ConvolutionLL(all, bll2)) != 0 {
		t.Errorf("Empty LL failed")
	}
}

func TestConvolutionMid(t *testing.T) {
	n, m := 1234, 2345
	a := make([]modint.Modint998244353, n)
	b := make([]modint.Modint998244353, m)
	for i := 0; i < n; i++ {
		a[i] = modint.NewModint998244353(rand.Int63())
	}
	for i := 0; i < m; i++ {
		b[i] = modint.NewModint998244353(rand.Int63())
	}

	expected := convNaive(a, b)
	actual := Convolution(a, b)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Mid failed")
	}
}

func TestConvolutionSimple(t *testing.T) {
	for n := 1; n < 20; n++ {
		for m := 1; m < 20; m++ {
			a := make([]modint.Modint998244353, n)
			b := make([]modint.Modint998244353, m)
			for i := 0; i < n; i++ {
				a[i] = modint.NewModint998244353(rand.Int63())
			}
			for i := 0; i < m; i++ {
				b[i] = modint.NewModint998244353(rand.Int63())
			}
			expected := convNaive(a, b)
			actual := Convolution(a, b)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Simple Mod998244353 failed for n=%d, m=%d", n, m)
			}
		}
	}
}

func TestConvolutionLL(t *testing.T) {
	for n := 1; n < 20; n++ {
		for m := 1; m < 20; m++ {
			a := make([]int64, n)
			b := make([]int64, m)
			for i := 0; i < n; i++ {
				a[i] = rand.Int63n(1000000) - 500000
			}
			for i := 0; i < m; i++ {
				b[i] = rand.Int63n(1000000) - 500000
			}
			expected := convLLNaive(a, b)
			actual := ConvolutionLL(a, b)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("ConvLL failed for n=%d, m=%d", n, m)
			}
		}
	}
}

func TestConvolutionLLBound(t *testing.T) {
	const MOD1 = uint64(469762049)
	const MOD2 = uint64(167772161)
	const MOD3 = uint64(754974721)
	M2M3 := MOD2 * MOD3
	M1M3 := MOD1 * MOD3
	M1M2 := MOD1 * MOD2

	for i := int64(-1000); i <= 1000; i++ {
		val := uint64(0)
		val -= M1M2
		val -= M1M3
		val -= M2M3
		val += uint64(i)
		a := []int64{int64(val)}
		b := []int64{1}
		if !reflect.DeepEqual(a, ConvolutionLL(a, b)) {
			t.Errorf("Bound failed for i=%d", i)
		}
	}

	for i := int64(0); i < 1000; i++ {
		a := []int64{math.MinInt64 + i}
		b := []int64{1}
		if !reflect.DeepEqual(a, ConvolutionLL(a, b)) {
			t.Errorf("Bound MinInt64 failed for i=%d", i)
		}
	}

	for i := int64(0); i < 1000; i++ {
		a := []int64{math.MaxInt64 - i}
		b := []int64{1}
		if !reflect.DeepEqual(a, ConvolutionLL(a, b)) {
			t.Errorf("Bound MaxInt64 failed for i=%d", i)
		}
	}
}

type Mod641 struct{}
func (Mod641) Mod() uint32 { return 641 }

func TestConvolution641(t *testing.T) {
	n, m := 64, 65
	a := make([]modint.StaticModint[Mod641], n)
	b := make([]modint.StaticModint[Mod641], m)
	for i := 0; i < n; i++ {
		a[i] = modint.NewStatic[Mod641](rand.Int63n(641))
	}
	for i := 0; i < m; i++ {
		b[i] = modint.NewStatic[Mod641](rand.Int63n(641))
	}
	if !reflect.DeepEqual(convNaive(a, b), Convolution(a, b)) {
		t.Errorf("Conv641 failed")
	}
}

type Mod2147483647 struct{}
func (Mod2147483647) Mod() uint32 { return 2147483647 }

func TestConvolution2147483647(t *testing.T) {
	a := []modint.StaticModint[Mod2147483647]{modint.NewStatic[Mod2147483647](rand.Int63n(2147483647))}
	b := []modint.StaticModint[Mod2147483647]{
		modint.NewStatic[Mod2147483647](rand.Int63n(2147483647)),
		modint.NewStatic[Mod2147483647](rand.Int63n(2147483647)),
	}
	if !reflect.DeepEqual(convNaive(a, b), Convolution(a, b)) {
		t.Errorf("Conv2147483647 failed")
	}
}

func TestConvolution18433(t *testing.T) {
	const mod = 18433
	type mint = modint.StaticModint[modint.Mod18433]
	a := make([]mint, 1)
	b := make([]mint, 1)
	a[0] = modint.NewStatic[modint.Mod18433](1)
	b[0] = modint.NewStatic[modint.Mod18433](2)
	c := Convolution(a, b)
	if len(c) != 1 || c[0].Val() != 2 {
		t.Errorf("Convolution18433 failed")
	}
}

func TestConvolution2(t *testing.T) {
	const mod = 2
	type mint = modint.StaticModint[modint.Mod2]
	a := make([]mint, 1)
	b := make([]mint, 1)
	a[0] = modint.NewStatic[modint.Mod2](1)
	b[0] = modint.NewStatic[modint.Mod2](1)
	c := Convolution(a, b)
	if len(c) != 1 || c[0].Val() != 1 {
		t.Errorf("Convolution2 failed")
	}
}

func TestConvolution257(t *testing.T) {
	const mod = 257
	type mint = modint.StaticModint[modint.Mod257]
	a := make([]mint, 1)
	b := make([]mint, 1)
	a[0] = modint.NewStatic[modint.Mod257](1)
	b[0] = modint.NewStatic[modint.Mod257](2)
	c := Convolution(a, b)
	if len(c) != 1 || c[0].Val() != 2 {
		t.Errorf("Convolution257 failed")
	}
}
