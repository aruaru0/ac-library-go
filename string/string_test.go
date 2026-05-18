package string

import (
	"reflect"
	"testing"
)

func lcpNaive(s, sa []int) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}
	lcp := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		l, r := sa[i], sa[i+1]
		for l+lcp[i] < n && r+lcp[i] < n && s[l+lcp[i]] == s[r+lcp[i]] {
			lcp[i]++
		}
	}
	return lcp
}

func zNaive(s []int) []int {
	n := len(s)
	z := make([]int, n)
	for i := 0; i < n; i++ {
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
	}
	return z
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestStringEmpty(t *testing.T) {
	if len(SuffixArrayString("")) != 0 {
		t.Errorf("SuffixArrayString empty failed")
	}
	if len(SuffixArray([]int{})) != 0 {
		t.Errorf("SuffixArray empty failed")
	}
	if len(ZAlgorithmString("")) != 0 {
		t.Errorf("ZAlgorithmString empty failed")
	}
	if len(ZAlgorithm([]int{})) != 0 {
		t.Errorf("ZAlgorithm empty failed")
	}
}

func TestStringSALCPNaive(t *testing.T) {
	for n := 1; n <= 5; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 4
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			maxC := 0
			for i := 0; i < n; i++ {
				s[i] = g % 4
				maxC = maxInt(maxC, s[i])
				g /= 4
			}
			sa := saNaive(s)
			if !reflect.DeepEqual(sa, SuffixArray(s)) {
				t.Errorf("SuffixArray mismatch")
			}
			if !reflect.DeepEqual(sa, SuffixArrayUpperBound(s, maxC)) {
				t.Errorf("SuffixArrayUpperBound mismatch")
			}
			if !reflect.DeepEqual(lcpNaive(s, sa), LcpArray(s, sa)) {
				t.Errorf("LcpArray mismatch")
			}
		}
	}
	for n := 1; n <= 10; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 2
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			maxC := 0
			for i := 0; i < n; i++ {
				s[i] = g % 2
				maxC = maxInt(maxC, s[i])
				g /= 2
			}
			sa := saNaive(s)
			if !reflect.DeepEqual(sa, SuffixArray(s)) {
				t.Errorf("SuffixArray mismatch")
			}
			if !reflect.DeepEqual(sa, SuffixArrayUpperBound(s, maxC)) {
				t.Errorf("SuffixArrayUpperBound mismatch")
			}
			if !reflect.DeepEqual(lcpNaive(s, sa), LcpArray(s, sa)) {
				t.Errorf("LcpArray mismatch")
			}
		}
	}
}

func TestStringInternalSANaiveNaive(t *testing.T) {
	for n := 1; n <= 5; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 4
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			for i := 0; i < n; i++ {
				s[i] = g % 4
				g /= 4
			}
			sa := saNaive(s)
			if !reflect.DeepEqual(saNaive(s), sa) {
				t.Errorf("saNaive mismatch")
			}
		}
	}
}

func TestStringInternalSADoublingNaive(t *testing.T) {
	for n := 1; n <= 5; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 4
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			for i := 0; i < n; i++ {
				s[i] = g % 4
				g /= 4
			}
			sa := saDoubling(s)
			if !reflect.DeepEqual(saNaive(s), sa) {
				t.Errorf("saDoubling mismatch")
			}
		}
	}
}

func TestStringInternalSAISNaive(t *testing.T) {
	for n := 1; n <= 5; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 4
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			maxC := 0
			for i := 0; i < n; i++ {
				s[i] = g % 4
				maxC = maxInt(maxC, s[i])
				g /= 4
			}
			sa := saIs(s, maxC)
			if !reflect.DeepEqual(saNaive(s), sa) {
				t.Errorf("saIs mismatch")
			}
		}
	}
}

func TestStringSAAllATest(t *testing.T) {
	for n := 1; n <= 100; n++ {
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = 10
		}
		saN := saNaive(s)
		if !reflect.DeepEqual(saN, SuffixArray(s)) {
			t.Errorf("SAAllA failed")
		}
		if !reflect.DeepEqual(saN, SuffixArrayUpperBound(s, 10)) {
			t.Errorf("SAAllA upperBound 10 failed")
		}
		if !reflect.DeepEqual(saN, SuffixArrayUpperBound(s, 12)) {
			t.Errorf("SAAllA upperBound 12 failed")
		}
	}
}

func TestStringSAAllABTest(t *testing.T) {
	for n := 1; n <= 100; n++ {
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = i % 2
		}
		saN := saNaive(s)
		if !reflect.DeepEqual(saN, SuffixArray(s)) {
			t.Errorf("SAAllAB failed")
		}
		if !reflect.DeepEqual(saN, SuffixArrayUpperBound(s, 3)) {
			t.Errorf("SAAllAB upperBound 3 failed")
		}
	}
	for n := 1; n <= 100; n++ {
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = 1 - (i % 2)
		}
		saN := saNaive(s)
		if !reflect.DeepEqual(saN, SuffixArray(s)) {
			t.Errorf("SAAllAB inv failed")
		}
		if !reflect.DeepEqual(saN, SuffixArrayUpperBound(s, 3)) {
			t.Errorf("SAAllAB inv upperBound 3 failed")
		}
	}
}

func TestStringSA(t *testing.T) {
	s := "missisippi"
	sa := SuffixArrayString(s)

	answer := []string{
		"i",
		"ippi",
		"isippi",
		"issisippi",
		"missisippi",
		"pi",
		"ppi",
		"sippi",
		"sisippi",
		"ssisippi",
	}

	if len(answer) != len(sa) {
		t.Errorf("SA length mismatch")
	}
	for i := 0; i < len(sa); i++ {
		if answer[i] != s[sa[i]:] {
			t.Errorf("SA content mismatch at %d", i)
		}
	}
}

func TestStringSASingle(t *testing.T) {
	if !reflect.DeepEqual([]int{0}, SuffixArray([]int{0})) {
		t.Errorf("SASingle 0 failed")
	}
	if !reflect.DeepEqual([]int{0}, SuffixArray([]int{-1})) {
		t.Errorf("SASingle -1 failed")
	}
	if !reflect.DeepEqual([]int{0}, SuffixArray([]int{1})) {
		t.Errorf("SASingle 1 failed")
	}
}

func TestStringLCP(t *testing.T) {
	s := "aab"
	sa := SuffixArrayString(s)
	if !reflect.DeepEqual([]int{0, 1, 2}, sa) {
		t.Errorf("SA aab failed")
	}
	lcp := LcpArrayString(s, sa)
	if !reflect.DeepEqual([]int{1, 0}, lcp) {
		t.Errorf("LCP aab failed")
	}

	if !reflect.DeepEqual(lcp, LcpArray([]int{0, 0, 1}, sa)) {
		t.Errorf("LCP 001 failed")
	}
	if !reflect.DeepEqual(lcp, LcpArray([]int{-100, -100, 100}, sa)) {
		t.Errorf("LCP -100 failed")
	}
}

func TestStringZAlgo(t *testing.T) {
	s := "abab"
	z := ZAlgorithmString(s)
	if !reflect.DeepEqual([]int{4, 0, 2, 0}, z) {
		t.Errorf("ZAlgo abab failed")
	}
	if !reflect.DeepEqual([]int{4, 0, 2, 0}, ZAlgorithm([]int{1, 10, 1, 10})) {
		t.Errorf("ZAlgo 1, 10, 1, 10 failed")
	}
	if !reflect.DeepEqual(zNaive([]int{0, 0, 0, 0, 0, 0, 0}), ZAlgorithm([]int{0, 0, 0, 0, 0, 0, 0})) {
		t.Errorf("ZAlgo all 0s failed")
	}
}

func TestStringZNaive(t *testing.T) {
	for n := 1; n <= 6; n++ {
		m := 1
		for i := 0; i < n; i++ {
			m *= 4
		}
		for f := 0; f < m; f++ {
			s := make([]int, n)
			g := f
			for i := 0; i < n; i++ {
				s[i] = g % 4
				g /= 4
			}
			if !reflect.DeepEqual(zNaive(s), ZAlgorithm(s)) {
				t.Errorf("ZNaive mismatch")
			}
		}
	}
}
