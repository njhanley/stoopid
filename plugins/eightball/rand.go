package eightball

import (
	"math/rand"
	"sync"
)

type rng struct {
	// immutable
	alias []int
	prob  []float64

	mu   sync.Mutex
	rand *rand.Rand
}

func newRNG(seed int64, p []float64) *rng {
	var s, l int
	small, large := make([]int, 0, len(p)), make([]int, 0, len(p))
	for i := range p {
		p[i] *= float64(len(p))
		if p[i] < 1 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	r := rng{
		alias: make([]int, len(p)),
		prob:  make([]float64, len(p)),
		rand:  rand.New(rand.NewSource(seed)),
	}

	for len(small) > 0 && len(large) > 0 {
		s, small = small[0], small[1:]
		l, large = large[0], large[1:]
		r.prob[s], r.alias[s] = p[s], l
		p[l] += p[s] - 1
		if p[l] < 1 {
			small = append(small, l)
		} else {
			large = append(large, l)
		}
	}
	for len(large) > 0 {
		l, large = large[0], large[1:]
		r.prob[l] = 1
	}
	for len(small) > 0 {
		s, small = small[0], small[1:]
		r.prob[s] = 1
	}

	return &r
}

func (r *rng) get() int {
	r.mu.Lock()
	x := r.rand.Float64() * float64(len(r.prob))
	r.mu.Unlock()

	i := int(x)
	if x-float64(i) < r.prob[i] {
		return i
	}
	return r.alias[i]
}
