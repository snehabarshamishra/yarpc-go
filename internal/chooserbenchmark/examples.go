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
