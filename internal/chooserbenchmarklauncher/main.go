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

	lbbench "go.uber.org/yarpc/internal/chooserbenchmark"
)

func main() {
	//config := &lbbench.Config{
	//	ClientGroups: []lbbench.ClientGroup{
	//		{
	//			Count:           1000,
	//			RPS:             1000,
	//			ListType:        lbbench.FewestPending,
	//			ListUpdaterType: lbbench.Static,
	//		},
	//	},
	//	ServerGroups: []lbbench.ServerGroup{
	//		{
	//			Name:          "fast",
	//			Count:         8,
	//			LatencyConfig: lbbench.RPSLatency(2000000),
	//		},
	//		{
	//			Name:          "slow",
	//			Count:         2,
	//			LatencyConfig: lbbench.RPSLatency(5),
	//		},
	//	},
	//	Duration: 3 * time.Second,
	//}
	//if err := lbbench.Run(config); err != nil {
	//	panic(err)
	//}
	config := &lbbench.Config{
		ClientGroups: []lbbench.ClientGroup{
			{
				Count:           1,
				RPS:             100,
				ListType:        lbbench.FewestPending,
				ListUpdaterType: lbbench.Static,
			},
		},
		ServerGroups: []lbbench.ServerGroup{
			{
				Name:          "slow",
				Count:         2,
				LatencyConfig: lbbench.RPSLatency(50),
			},
		},
		Duration: 3 * time.Second,
	}
	if err := lbbench.Run(config); err != nil {
		panic(err)
	}
}
