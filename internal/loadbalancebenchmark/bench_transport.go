package loadbalancingbenchmark

import (
	"strconv"

	"go.uber.org/yarpc/api/peer"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/pkg/lifecycletest"
)

type BenchTransport struct {
	transport.Lifecycle
}

func NewBenchTransport() *BenchTransport {
	t := &BenchTransport{
		Lifecycle: lifecycletest.NewNop(),
	}
	return t
}

func (t *BenchTransport) RetainPeer(id peer.Identifier, ps peer.Subscriber) (peer.Peer, error) {
	i, err := strconv.Atoi(id.Identifier())
	if err != nil {
		return nil, err
	}
	return NewBenchPeer(i), nil
}

// TODO update release peer logic if we want to simulate server break down and come back
func (t *BenchTransport) ReleasePeer(id peer.Identifier, ps peer.Subscriber) error {
	return nil
}
