package loadbalancingbenchmark

import (
	"strconv"
)

type BenchIdentifier struct {
	id int
}

func (p BenchIdentifier) Identifier() string {
	return strconv.Itoa(p.id)
}
