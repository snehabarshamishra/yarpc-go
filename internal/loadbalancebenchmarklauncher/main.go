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

package main

import (
	"time"

	lbbench "go.uber.org/yarpc/internal/loadbalancebenchmark"
)

func main() {
	config := &lbbench.TestConfig{
		ClientGroup: []lbbench.ClientGroup{
			{
				Cnt:    1000,
				Rps:    1000,
				LType:  lbbench.FewestPending,
				LUType: lbbench.Static,
			},
		},
		ServerGroup: []lbbench.ServerGroup{
			{
				Cnt:           8,
				Type:          lbbench.NormalMachine,
				LatencyConfig: lbbench.RpsLatency(2000000),
			},
			{
				Cnt:           2,
				Type:          lbbench.SlowMachine,
				LatencyConfig: lbbench.RpsLatency(5),
			},
		},
		Duration: 3 * time.Second,
	}
	if err := lbbench.StartBenchmark(config); err != nil {
		panic(err)
	}
}
