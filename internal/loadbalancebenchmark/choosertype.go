package loadbalancingbenchmark

type ChooserType string

const (
	Unknown       = "unknown"
	RoundRobin    = "roundrobin"
	FewestPending = "fewestpending"
)
