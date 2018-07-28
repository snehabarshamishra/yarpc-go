// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package chooserbenchmark

import (
	"math"
	"math/rand"
	"time"
)

// LogNormalSigma is a hard-coded parameter we find suitable to simulate real
// world latency in range of [1ms, 10s]
const LogNormalSigma = 0.5

// LogNormalLatency determines the duration of sleep time, we use log normal
// distribution to simulate latency, http://www.lognormal.com/features/. log
// normal means a random variable whose log is normally distributed, formulas
// are referenced from https://en.wikipedia.org/wiki/Log-normal_distribution
type LogNormalLatency struct {
	mu     float64
	sigma  float64
	median time.Duration
}

// NewLogNormalLatency returns a log normal distribution random generator that
// takes the input latency as median
func NewLogNormalLatency(latency time.Duration) *LogNormalLatency {
	median := latency
	mu := math.Log(float64(median))

	return &LogNormalLatency{
		mu:     mu,
		sigma:  LogNormalSigma,
		median: median,
	}
}

// Random returns a service delay obey to log normal distribution
func (l *LogNormalLatency) Random() time.Duration {
	rnd := rand.NormFloat64()
	return time.Duration(math.Exp(rnd*l.sigma + l.mu))
}

// CDF a.k.a. Cumulative Density Function, return the probability that Random()
// takes a value smaller than or equal to x.
// you can use CDF to calculate p99, p90, p50 etc, using binary search
func (l *LogNormalLatency) CDF(x float64) float64 {
	return 0.5 + 0.5*math.Erf((math.Log(x)-l.mu)/(math.Sqrt2*l.sigma))
}
