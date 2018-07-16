package main

import (
	"time"

	lbbench "go.uber.org/yarpc/internal/loadbalancebenchmark"
)

func main() {
	config := &lbbench.TestConfig{
		ClientGroup: []lbbench.ClientGroup{
			{
				Cnt:    100,
				Rps:    100,
				LType:  lbbench.FewestPending,
				LUType: lbbench.Static,
			},
		},
		ServerGroup: []lbbench.ServerGroup{
			{
				Cnt:           8,
				Type:          lbbench.NormalMachine,
				LatencyConfig: lbbench.RpsLatency(100),
			},
			{
				Cnt:           1,
				Type:          lbbench.SlowMachine,
				LatencyConfig: lbbench.RpsLatency(1),
			},
		},
		Duration: 3 * time.Second,
	}
	if err := lbbench.StartBenchmark(config); err != nil {
		panic(err)
	}
}
