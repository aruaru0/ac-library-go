package string

import (
	"sort"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func saNaive(s []int) []int {
	n := len(s)
	sa := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
	}
	sort.Slice(sa, func(i, j int) bool {
		l, r := sa[i], sa[j]
		if l == r {
			return false
		}
		for l < n && r < n {
			if s[l] != s[r] {
				return s[l] < s[r]
			}
			l++
			r++
		}
		return l == n
	})
	return sa
}

func saDoubling(s []int) []int {
	n := len(s)
	sa := make([]int, n)
	rnk := make([]int, n)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
		rnk[i] = s[i]
	}

	for k := 1; k < n; k *= 2 {
		cmp := func(x, y int) bool {
			if rnk[x] != rnk[y] {
				return rnk[x] < rnk[y]
			}
			rx, ry := -1, -1
			if x+k < n {
				rx = rnk[x+k]
			}
			if y+k < n {
				ry = rnk[y+k]
			}
			return rx < ry
		}
		sort.Slice(sa, func(i, j int) bool {
			return cmp(sa[i], sa[j])
		})
		tmp[sa[0]] = 0
		for i := 1; i < n; i++ {
			tmp[sa[i]] = tmp[sa[i-1]]
			if cmp(sa[i-1], sa[i]) {
				tmp[sa[i]]++
			}
		}
		tmp, rnk = rnk, tmp
	}
	return sa
}

const thresholdNaive = 10
const thresholdDoubling = 40

func saIs(s []int, upper int) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}
	if n == 1 {
		return []int{0}
	}
	if n == 2 {
		if s[0] < s[1] {
			return []int{0, 1}
		} else {
			return []int{1, 0}
		}
	}
	if n < thresholdNaive {
		return saNaive(s)
	}
	if n < thresholdDoubling {
		return saDoubling(s)
	}

	sa := make([]int, n)
	ls := make([]bool, n)
	for i := n - 2; i >= 0; i-- {
		if s[i] == s[i+1] {
			ls[i] = ls[i+1]
		} else {
			ls[i] = s[i] < s[i+1]
		}
	}

	sumL := make([]int, upper+1)
	sumS := make([]int, upper+1)
	for i := 0; i < n; i++ {
		if !ls[i] {
			sumS[s[i]]++
		} else {
			sumL[s[i]+1]++
		}
	}
	for i := 0; i <= upper; i++ {
		sumS[i] += sumL[i]
		if i < upper {
			sumL[i+1] += sumS[i]
		}
	}

	induce := func(lms []int) {
		for i := range sa {
			sa[i] = -1
		}
		buf := make([]int, upper+1)
		copy(buf, sumS)
		for _, d := range lms {
			if d == n {
				continue
			}
			sa[buf[s[d]]] = d
			buf[s[d]]++
		}
		copy(buf, sumL)
		sa[buf[s[n-1]]] = n - 1
		buf[s[n-1]]++
		for i := 0; i < n; i++ {
			v := sa[i]
			if v >= 1 && !ls[v-1] {
				sa[buf[s[v-1]]] = v - 1
				buf[s[v-1]]++
			}
		}
		copy(buf, sumL)
		for i := n - 1; i >= 0; i-- {
			v := sa[i]
			if v >= 1 && ls[v-1] {
				buf[s[v-1]+1]--
				sa[buf[s[v-1]+1]] = v - 1
			}
		}
	}

	lmsMap := make([]int, n+1)
	for i := range lmsMap {
		lmsMap[i] = -1
	}
	m := 0
	for i := 1; i < n; i++ {
		if !ls[i-1] && ls[i] {
			lmsMap[i] = m
			m++
		}
	}
	lms := make([]int, 0, m)
	for i := 1; i < n; i++ {
		if !ls[i-1] && ls[i] {
			lms = append(lms, i)
		}
	}

	induce(lms)

	if m > 0 {
		sortedLms := make([]int, 0, m)
		for _, v := range sa {
			if lmsMap[v] != -1 {
				sortedLms = append(sortedLms, v)
			}
		}
		recS := make([]int, m)
		recUpper := 0
		recS[lmsMap[sortedLms[0]]] = 0
		for i := 1; i < m; i++ {
			l := sortedLms[i-1]
			r := sortedLms[i]
			endL := n
			if lmsMap[l]+1 < m {
				endL = lms[lmsMap[l]+1]
			}
			endR := n
			if lmsMap[r]+1 < m {
				endR = lms[lmsMap[r]+1]
			}
			same := true
			if endL-l != endR-r {
				same = false
			} else {
				for l < endL {
					if s[l] != s[r] {
						break
					}
					l++
					r++
				}
				if l == n || s[l] != s[r] {
					same = false
				}
			}
			if !same {
				recUpper++
			}
			recS[lmsMap[sortedLms[i]]] = recUpper
		}

		recSa := saIs(recS, recUpper)

		for i := 0; i < m; i++ {
			sortedLms[i] = lms[recSa[i]]
		}
		induce(sortedLms)
	}
	return sa
}

// SuffixArrayUpperBound calculates the suffix array of s where elements are in [0, upper].
func SuffixArrayUpperBound(s []int, upper int) []int {
	if upper < 0 {
		panic("upper must be non-negative")
	}
	for _, d := range s {
		if d < 0 || d > upper {
			panic("element out of bounds")
		}
	}
	return saIs(s, upper)
}

// SuffixArray calculates the suffix array of s.
func SuffixArray[T Integer](s []T) []int {
	n := len(s)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		return s[idx[i]] < s[idx[j]]
	})
	s2 := make([]int, n)
	now := 0
	for i := 0; i < n; i++ {
		if i > 0 && s[idx[i-1]] != s[idx[i]] {
			now++
		}
		s2[idx[i]] = now
	}
	return saIs(s2, now)
}

// SuffixArrayString calculates the suffix array of a string.
func SuffixArrayString(s string) []int {
	n := len(s)
	s2 := make([]int, n)
	for i := 0; i < n; i++ {
		s2[i] = int(s[i])
	}
	return saIs(s2, 255)
}

// LcpArray calculates the LCP array from the suffix array and the original slice.
func LcpArray[T Integer](s []T, sa []int) []int {
	n := len(s)
	if n != len(sa) {
		panic("length of s and sa must be equal")
	}
	if n < 1 {
		panic("n must be at least 1")
	}
	rnk := make([]int, n)
	for i := 0; i < n; i++ {
		if sa[i] < 0 || sa[i] >= n {
			panic("sa contains out of bounds elements")
		}
		rnk[sa[i]] = i
	}
	lcp := make([]int, n-1)
	h := 0
	for i := 0; i < n; i++ {
		if h > 0 {
			h--
		}
		if rnk[i] == 0 {
			continue
		}
		j := sa[rnk[i]-1]
		for j+h < n && i+h < n {
			if s[j+h] != s[i+h] {
				break
			}
			h++
		}
		lcp[rnk[i]-1] = h
	}
	return lcp
}

// LcpArrayString calculates the LCP array from the suffix array and the original string.
func LcpArrayString(s string, sa []int) []int {
	n := len(s)
	s2 := make([]int, n)
	for i := 0; i < n; i++ {
		s2[i] = int(s[i])
	}
	return LcpArray(s2, sa)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ZAlgorithm calculates the Z array of a slice.
func ZAlgorithm[T Integer](s []T) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}
	z := make([]int, n)
	z[0] = 0
	j := 0
	for i := 1; i < n; i++ {
		k := 0
		if j+z[j] > i {
			k = minInt(j+z[j]-i, z[i-j])
		}
		for i+k < n && s[k] == s[i+k] {
			k++
		}
		z[i] = k
		if j+z[j] < i+z[i] {
			j = i
		}
	}
	z[0] = n
	return z
}

// ZAlgorithmString calculates the Z array of a string.
func ZAlgorithmString(s string) []int {
	n := len(s)
	s2 := make([]int, n)
	for i := 0; i < n; i++ {
		s2[i] = int(s[i])
	}
	return ZAlgorithm(s2)
}
