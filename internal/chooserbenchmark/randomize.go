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

// LogNormalLatency determines the duration of sleep time on server side, we
// use log normal distribution to simulate latency. log normal means a random
// variable whose log is normally distributed, formulas are referenced
// from https://en.wikipedia.org/wiki/Log-normal_distribution
type LogNormalLatency struct {
	mu     float64
	sigma  float64
	median float64
}

// NewLogNormalLatency returns a log normal distribution random generator that
// takes the input latency as median
func NewLogNormalLatency(latency time.Duration) *LogNormalLatency {
	median := float64(latency)
	mu := math.Log(median)

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

// Median returns the median of sleep time on server side
func (l *LogNormalLatency) Median() time.Duration {
	return time.Duration(l.median)
}

// CDF a.k.a. Cumulative Density Function, return the probability that Random()
// takes a value smaller than or equal to x.
// you can use CDF to calculate p99, p90, p50 etc, using binary search
func (l *LogNormalLatency) CDF(x float64) float64 {
	return 0.5 + 0.5*math.Erf((math.Log(x)-l.mu)/(math.Sqrt2*l.sigma))
}

// NormalDistSleepTime determines the duration of sleep time on client side, we
// use normal distribution to increase randomness
type NormalDistSleepTime struct {
	mu    float64
	sigma float64
}

// NewNormalDistSleepTime returns a normal distribution random generator that
// takes the mu, sigma calculated from client RPS in config
func NewNormalDistSleepTime(rps int) *NormalDistSleepTime {
	sleepTime := float64(time.Second) / float64(rps)
	return &NormalDistSleepTime{
		mu:    sleepTime,
		sigma: sleepTime / 20, // the deviation is used to increase randomness
	}
}

// Random returns a client sleep time obey to normal distribution
func (n *NormalDistSleepTime) Random() time.Duration {
	return time.Duration(rand.NormFloat64()*n.sigma + n.mu)
}

// Median returns the median of sleep time on client side
func (n *NormalDistSleepTime) Median() time.Duration {
	return time.Duration(n.mu)
}
