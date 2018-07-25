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
	"time"
)

var RoundRobinWorks = &Config{
	ClientGroups: []ClientGroup{
		{
			Name:        "roundrobin",
			Count:       500,
			RPS:         20,
			Constructor: RoundRobin,
		},
		{
			Name:        "pendingheap",
			Count:       500,
			RPS:         20,
			Constructor: PendingHeap,
		},
	},
	ServerGroups: []ServerGroup{
		{
			Name:          "normal",
			Count:         50,
			LatencyConfig: time.Millisecond * 100,
		},
	},
	Duration: 10 * time.Second,
}

var FewestPendingSuperior = &Config{
	ClientGroups: []ClientGroup{
		{
			Name:        "roundrobin",
			Count:       1000,
			RPS:         20,
			Constructor: RoundRobin,
		},
		{
			Name:        "pendingheap",
			Count:       1000,
			RPS:         20,
			Constructor: PendingHeap,
		},
	},
	ServerGroups: []ServerGroup{
		{
			Name:          "normal",
			Count:         5,
			LatencyConfig: time.Millisecond * 100,
		},
		{
			Name:          "slow",
			Count:         5,
			LatencyConfig: time.Second,
		},
	},
	Duration: 10 * time.Second,
}

var FewestPendingDegradation = &Config{
	ClientGroups: []ClientGroup{
		{
			Name:        "roundrobin",
			Count:       1000,
			RPS:         20,
			Constructor: RoundRobin,
		},
		{
			Name:        "pendingheap",
			Count:       1000,
			RPS:         20,
			Constructor: PendingHeap,
		},
	},
	ServerGroups: []ServerGroup{
		{
			Name:          "normal",
			Count:         50,
			LatencyConfig: time.Millisecond * 100,
		},
		{
			Name:          "slow",
			Count:         50,
			LatencyConfig: time.Second,
		},
	},
	Duration: 10 * time.Second,
}
