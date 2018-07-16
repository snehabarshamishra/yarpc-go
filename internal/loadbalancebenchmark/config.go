package loadbalancingbenchmark

import (
	"fmt"
	"time"
)

type ClientGroup struct {
	Cnt    int
	Rps    int
	LType  ListType
	LUType ListUpdaterType
}

type LatencyConfig struct {
	P50  time.Duration
	P90  time.Duration
	P99  time.Duration
	P100 time.Duration
}

type ServerGroup struct {
	Cnt           int
	Type          MachineType
	LatencyConfig *LatencyConfig
}

type TestConfig struct {
	ClientGroup []ClientGroup
	ServerGroup []ServerGroup
	Duration    time.Duration
}

// return a normal latency config that satisfy the given rps
func RpsLatency(rps int) *LatencyConfig {
	// measure in millisecond
	normal := time.Second / time.Duration(rps)
	return &LatencyConfig{
		P50:  time.Duration(int(0.4 * float32(normal))),
		P90:  time.Duration(int(0.7 * float32(normal))),
		P99:  time.Duration(int(0.8 * float32(normal))),
		P100: time.Duration(int(2.0 * float32(normal))),
	}
}

func (config *TestConfig) GetClientTotalCount() int {
	total := 0
	for _, group := range config.ClientGroup {
		total += group.Cnt
	}
	return total
}

func (config *TestConfig) GetServerTotalCount() int {
	total := 0
	for _, group := range config.ServerGroup {
		total += group.Cnt
	}
	return total
}

func (clientGroup ClientGroup) checkListType() error {
	listType := clientGroup.LType
	if _, ok := SupportedListType[listType]; !ok {
		return fmt.Errorf(`list type %q is not supported`, listType)
	}
	return nil
}

func (clientGroup ClientGroup) checkListUpdaterType() error {
	listUpdaterType := clientGroup.LUType
	if _, ok := SupportedListUpdaterType[listUpdaterType]; !ok {
		return fmt.Errorf(`list updater type %q is not supported`, listUpdaterType)
	}
	return nil
}

func (config *TestConfig) checkClientGroup() error {
	clientGroup := config.ClientGroup
	for _, group := range clientGroup {
		if err := group.checkListType(); err != nil {
			return err
		}
		if err := group.checkListUpdaterType(); err != nil {
			return err
		}
		if group.Rps < 0 {
			return fmt.Errorf("rps field should be greater than 0 rps: %d", group.Rps)
		}
	}
	return nil
}

func (config *TestConfig) checkServerGroup() error {
	serverGroup := config.ServerGroup
	for _, group := range serverGroup {
		if _, ok := SupportedMachineType[group.Type]; !ok {
			return fmt.Errorf(`machine type %q is not supported`, group.Type)
		}
		latencyConfig := group.LatencyConfig
		p50, p90, p99, p100 := latencyConfig.P50, latencyConfig.P90, latencyConfig.P99, latencyConfig.P100
		if p50 < 0 || p90 < p50 || p99 < p90 || p100 < p99 {
			return fmt.Errorf("latency profile inconsistent p50: %v, p90: %v, p99: %v, p100: %v", p50, p90, p99, p100)
		}
	}
	return nil
}

func (config *TestConfig) Validate() error {
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
