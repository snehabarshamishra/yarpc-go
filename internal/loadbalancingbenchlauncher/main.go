package main

import (
	lbbench "go.uber.org/yarpc/internal/loadbalancebenchmark"
)

func main() {
	config, _ := lbbench.NewTestConfig(map[string]string{
		"ChooserType": lbbench.RoundRobin,
		"ServerCount": "3",
	})
	lbbench.StartBenchmark(config)
}
