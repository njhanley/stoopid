package eightball

import (
	"time"
)

type eightball struct {
	answers []string
	rng     *rng
}

type answer struct {
	Text   string
	Weight float64
}

func newEightball(x []answer) *eightball {
	answers := make([]string, len(x))
	weights := make([]float64, len(x))
	var sum float64
	for i, v := range x {
		answers[i], weights[i] = v.Text, v.Weight
		sum += weights[i]
	}

	// normalize to [0,1]
	for i := range weights {
		weights[i] /= sum
	}

	return &eightball{
		answers: answers,
		rng:     newRNG(time.Now().UnixNano(), weights),
	}
}

func (e *eightball) answer() string {
	return e.answers[e.rng.get()]
}
