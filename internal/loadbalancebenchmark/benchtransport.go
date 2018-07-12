package loadbalancingbenchmark

import (
	"strconv"

	"go.uber.org/yarpc/api/peer"
)

type BenchTransport struct {
}

func NewBenchTransport() *BenchTransport {
	return &BenchTransport{}
}

func (t *BenchTransport) RetainPeer(id peer.Identifier, ps peer.Subscriber) (peer.Peer, error) {
	i, err := strconv.Atoi(id.Identifier())
	if err != nil {
		return nil, err
	}
	return NewBenchPeer(i, ps), nil
}

// TODO update release peer logic if we want to simulate server break down and come back
func (t *BenchTransport) ReleasePeer(id peer.Identifier, ps peer.Subscriber) error {
	return nil
}
