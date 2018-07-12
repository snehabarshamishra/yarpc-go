package loadbalancingbenchmark

type ListType string

const (
	RoundRobin    = "roundrobin"
	FewestPending = "fewestpending"
)

var SupportedListType = map[ListType]struct{}{
	RoundRobin:    {},
	FewestPending: {},
}

type ListUpdaterType string

const (
	Static                  = "static"
	RandomSubsetting        = "randomsubsetting"
	DeterministicSubsetting = "deterministicsubsetting"
)

var SupportedListUpdaterType = map[ListUpdaterType]struct{}{
	Static:                  {},
	RandomSubsetting:        {},
	DeterministicSubsetting: {},
}

type MachineType string

const (
	CustomMachine = "custom"
	SlowMachine   = "slow"
	NormalMachine = "normal"
	FastMachine   = "fast"
)

var SupportedMachineType = map[MachineType]struct{}{
	CustomMachine: {},
	SlowMachine:   {},
	NormalMachine: {},
	FastMachine:   {},
}
