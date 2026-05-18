package lazysegtree

import (
	"math/rand"
	"testing"
)

type timeManager struct {
	v []int
}

func newTimeManager(n int) *timeManager {
	v := make([]int, n)
	for i := range v {
		v[i] = -1
	}
	return &timeManager{v: v}
}

func (tm *timeManager) action(l, r, time int) {
	for i := l; i < r; i++ {
		tm.v[i] = time
	}
}

func (tm *timeManager) prod(l, r int) int {
	res := -1
	for i := l; i < r; i++ {
		if tm.v[i] > res {
			res = tm.v[i]
		}
	}
	return res
}

type sStruct struct {
	l, r, time int
}

type tStruct struct {
	newTime int
}

func opSSStruct(l, r sStruct) sStruct {
	if l.l == -1 {
		return r
	}
	if r.l == -1 {
		return l
	}
	if l.r != r.l {
		panic("invalid l.r != r.l")
	}
	time := l.time
	if r.time > time {
		time = r.time
	}
	return sStruct{l.l, r.r, time}
}

func opTSStruct(l tStruct, r sStruct) sStruct {
	if l.newTime == -1 {
		return r
	}
	if r.time >= l.newTime {
		panic("invalid r.time >= l.newTime")
	}
	return sStruct{r.l, r.r, l.newTime}
}

func opTTStruct(l, r tStruct) tStruct {
	if l.newTime == -1 {
		return r
	}
	if r.newTime == -1 {
		return l
	}
	if l.newTime <= r.newTime {
		panic("invalid l.newTime <= r.newTime")
	}
	return l
}

func eSStruct() sStruct {
	return sStruct{-1, -1, -1}
}

func eTStruct() tStruct {
	return tStruct{-1}
}

func TestLazySegTreeStressNaive(t *testing.T) {
	for n := 1; n <= 30; n++ {
		for ph := 0; ph < 10; ph++ {
			seg0 := NewLazySegTree(n, opSSStruct, eSStruct, opTSStruct, opTTStruct, eTStruct)
			tm := newTimeManager(n)
			for i := 0; i < n; i++ {
				seg0.Set(i, sStruct{i, i + 1, -1})
			}
			now := 0
			for q := 0; q < 3000; q++ {
				ty := rand.Intn(4)
				l := rand.Intn(n + 1)
				r := rand.Intn(n + 1)
				if l > r {
					l, r = r, l
				}

				if ty == 0 {
					res := seg0.Prod(l, r)
					if l != r {
						if l != res.l || r != res.r || tm.prod(l, r) != res.time {
							t.Errorf("Prod mismatch: l=%d r=%d res=%+v tm=%d", l, r, res, tm.prod(l, r))
						}
					} else {
						if res.l != -1 {
							t.Errorf("Prod mismatch empty: res=%+v", res)
						}
					}
				} else if ty == 1 {
					if l == n {
						l-- // Ensure valid index
					}
					res := seg0.Get(l)
					if l != res.l || l+1 != res.r || tm.prod(l, l+1) != res.time {
						t.Errorf("Get mismatch: l=%d res=%+v tm=%d", l, res, tm.prod(l, l+1))
					}
				} else if ty == 2 {
					now++
					seg0.ApplyRange(l, r, tStruct{now})
					tm.action(l, r, now)
				} else if ty == 3 {
					if l == n {
						l--
					}
					now++
					seg0.Apply(l, tStruct{now})
					tm.action(l, l+1, now)
				}
			}
		}
	}
}

func TestLazySegTreeStressMaxRight(t *testing.T) {
	for n := 1; n <= 30; n++ {
		for ph := 0; ph < 10; ph++ {
			seg0 := NewLazySegTree(n, opSSStruct, eSStruct, opTSStruct, opTTStruct, eTStruct)
			tm := newTimeManager(n)
			for i := 0; i < n; i++ {
				seg0.Set(i, sStruct{i, i + 1, -1})
			}
			now := 0
			for q := 0; q < 1000; q++ {
				ty := rand.Intn(2)
				l := rand.Intn(n + 1)
				r := rand.Intn(n + 1)
				if l > r {
					l, r = r, l
				}

				if ty == 0 {
					res := seg0.MaxRight(l, func(s sStruct) bool {
						if s.l == -1 {
							return true
						}
						if s.l != l {
							panic("s.l != l")
						}
						if s.time != tm.prod(l, s.r) {
							panic("s.time mismatch")
						}
						return s.r <= r
					})
					if res != r {
						t.Errorf("MaxRight mismatch: l=%d r=%d res=%d", l, r, res)
					}
				} else {
					now++
					seg0.ApplyRange(l, r, tStruct{now})
					tm.action(l, r, now)
				}
			}
		}
	}
}

func TestLazySegTreeStressMinLeft(t *testing.T) {
	for n := 1; n <= 30; n++ {
		for ph := 0; ph < 10; ph++ {
			seg0 := NewLazySegTree(n, opSSStruct, eSStruct, opTSStruct, opTTStruct, eTStruct)
			tm := newTimeManager(n)
			for i := 0; i < n; i++ {
				seg0.Set(i, sStruct{i, i + 1, -1})
			}
			now := 0
			for q := 0; q < 1000; q++ {
				ty := rand.Intn(2)
				l := rand.Intn(n + 1)
				r := rand.Intn(n + 1)
				if l > r {
					l, r = r, l
				}

				if ty == 0 {
					res := seg0.MinLeft(r, func(s sStruct) bool {
						if s.l == -1 {
							return true
						}
						if s.r != r {
							panic("s.r != r")
						}
						if s.time != tm.prod(s.l, r) {
							panic("s.time mismatch")
						}
						return l <= s.l
					})
					if res != l {
						t.Errorf("MinLeft mismatch: l=%d r=%d res=%d", l, r, res)
					}
				} else {
					now++
					seg0.ApplyRange(l, r, tStruct{now})
					tm.action(l, r, now)
				}
			}
		}
	}
}
