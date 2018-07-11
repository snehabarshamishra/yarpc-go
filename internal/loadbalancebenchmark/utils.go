package loadbalancingbenchmark

import (
	"go.uber.org/yarpc/api/peer"
)

func MakePeerIdentifier(n int) []peer.Identifier {
	if n <= 0 {
		n = 0
	}
	ids := make([]peer.Identifier, n)
	for i := 0; i < n; i++ {
		ids[i] = BenchIdentifier{id: i}
	}
	return ids
}
