package eightball

import (
	"math/rand"
	"sync"
)

// Weighted random number generator
type wrng struct {
	mu   sync.Mutex
	rand *rand.Rand

	// immutable
	alias []int
	prob  []float64
}

func newWRNG(seed int64, prob []float64) *wrng {
	n := len(prob)
	small, large := make([]int, 0, n), make([]int, 0, n)
	for i := range prob {
		prob[i] *= float64(n)
		if prob[i] < 1 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	w := wrng{
		rand:  rand.New(rand.NewSource(seed)),
		alias: make([]int, n),
		prob:  make([]float64, n),
	}

	var s, l int
	for len(small) > 0 && len(large) > 0 {
		s, small = small[0], small[1:]
		l, large = large[0], large[1:]

		w.prob[s] = prob[s]
		w.alias[s] = l

		prob[l] += prob[s] - 1
		if prob[l] < 1 {
			small = append(small, l)
		} else {
			large = append(large, l)
		}
	}

	for len(large) > 0 {
		l, large = large[0], large[1:]
		w.prob[l] = 1
	}

	for len(small) > 0 {
		s, small = small[0], small[1:]
		w.prob[s] = 1
	}

	return &w
}

func (w *wrng) get() int {
	w.mu.Lock()
	x := w.rand.Float64() * float64(len(w.prob))
	w.mu.Unlock()

	i := int(x)
	if x-float64(i) < w.prob[i] {
		return i
	}
	return w.alias[i]
}
