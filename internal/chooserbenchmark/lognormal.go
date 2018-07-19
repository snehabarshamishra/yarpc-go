package chooserbenchmark

import (
	"math"
	"math/rand"
	"time"
)

type LogNormalLatency struct {
	mu     float64
	sigma  float64
	median time.Duration
	// pin is necessary to reduce overhead in large scale, 10 millions random
	// call will need 500ms, os will provide enough randomness in that scale.
	pin bool
}

func NewLogNormalLatency(latencyConfig *LatencyConfig) *LogNormalLatency {
	median, pin := latencyConfig.Median, latencyConfig.Pin
	mu := math.Log(float64(median))

	return &LogNormalLatency{
		mu:     mu,
		sigma:  LogNormalSigma,
		median: median,
		pin:    pin,
	}
}

func (l *LogNormalLatency) Random() time.Duration {
	if l.pin {
		return l.median
	}
	rnd := rand.NormFloat64()
	return time.Duration(math.Exp(rnd*l.sigma + l.mu))
}
