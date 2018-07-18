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
	"fmt"
	"time"
)

type Config struct {
	ClientGroups []ClientGroup
	ServerGroups []ServerGroup
	Duration     time.Duration

	CustomListTypes        []CustomListType
	CustomListUpdaterTypes []CustomListUpdaterType
}

type ClientGroup struct {
	Name            string
	Count           int
	RPS             int
	ListType        ListType
	ListUpdaterType ListUpdaterType
}

type ServerGroup struct {
	Name          string
	Count         int
	LatencyConfig *LatencyConfig
}

type CustomListType struct {
	ListType   ListType
	ListMethod ListFactoryMethod
}

type CustomListUpdaterType struct {
	ListUpdaterType ListUpdaterType
	UpdaterMethod   ListUpdaterFactoryMethod
}

type LatencyConfig struct {
	P50  time.Duration
	P90  time.Duration
	P99  time.Duration
	P100 time.Duration
}

// return a normal latency config that satisfy the given rps
func RPSLatency(rps int) *LatencyConfig {
	if rps <= 0 {
		panic("rps should be greater than 0")
	}
	// measure in millisecond
	normal := time.Second / time.Duration(rps)
	return &LatencyConfig{
		P50:  time.Duration(int(0.4 * float32(normal))),
		P90:  time.Duration(int(0.7 * float32(normal))),
		P99:  time.Duration(int(0.8 * float32(normal))),
		P100: time.Duration(int(2.0 * float32(normal))),
	}
}

func (config *Config) checkClientGroup() error {
	clientGroup := config.ClientGroups
	for _, group := range clientGroup {
		if group.RPS < 0 {
			return fmt.Errorf("rps field should be greater than 0 rps: %d", group.RPS)
		}
		if group.Count <= 0 {
			return fmt.Errorf("client group count must be greater than 0, client group count: %d", group.Count)
		}
	}
	return nil
}

func (config *Config) checkServerGroup() error {
	serverGroup := config.ServerGroups
	for _, group := range serverGroup {
		latencyConfig := group.LatencyConfig
		p50, p90, p99, p100 := latencyConfig.P50, latencyConfig.P90, latencyConfig.P99, latencyConfig.P100
		if p50 < 0 || p90 < p50 || p99 < p90 || p100 < p99 {
			return fmt.Errorf("latency profile inconsistent p50: %v, p90: %v, p99: %v, p100: %v", p50, p90, p99, p100)
		}
		if group.Count <= 0 {
			return fmt.Errorf("server group count must be greater than 0, server group count: %d", group.Count)
		}
	}
	return nil
}

func (config *Config) Validate() error {
	fmt.Println("checking config...")
	if config.Duration <= 0 {
		return fmt.Errorf(`test duration should be greater than 0, current: %v`, config.Duration)
	}
	if err := config.checkClientGroup(); err != nil {
		return err
	}
	if err := config.checkServerGroup(); err != nil {
		return err
	}
	return nil
}
