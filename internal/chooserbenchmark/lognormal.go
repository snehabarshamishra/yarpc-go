package chooserbenchmark

import (
	"math"
	"math/rand"
	"time"
)

const LogNormalSigma = 0.5

type LogNormalLatency struct {
	mu     float64
	sigma  float64
	median time.Duration
	// pin is necessary to reduce overhead in large scale, 10 millions random
	// call will need 500ms, os will provide enough randomness in that scale.
	pin bool
}

func NewLogNormalLatency(latency time.Duration) *LogNormalLatency {
	median := latency
	mu := math.Log(float64(median))

	return &LogNormalLatency{
		mu:     mu,
		sigma:  LogNormalSigma,
		median: median,
	}
}

func (l *LogNormalLatency) Random() time.Duration {
	rnd := rand.NormFloat64()
	return time.Duration(math.Exp(rnd*l.sigma + l.mu))
}

func (l *LogNormalLatency) CDF(x float64) float64 {
	return 0
}
