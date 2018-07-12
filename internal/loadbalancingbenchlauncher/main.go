package main

import (
	lbbench "go.uber.org/yarpc/internal/loadbalancebenchmark"
)

func main() {
	config, _ := lbbench.NewTestConfig(map[string]string{
		"ChooserType": lbbench.FewestPending,
		"ServerCount": "3",
		"ClientCount": "3",
	})
	if err := lbbench.StartBenchmark(config); err != nil {
		panic(err)
	}

}
