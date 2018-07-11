package loadbalancingbenchmark

import (
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/peer"
)

type BenchPeer struct {
	id      BenchIdentifier
	pending atomic.Int32
}

func (p *BenchPeer) Identifier() string {
	return p.id.Identifier()
}

func NewBenchPeer(id int) *BenchPeer {
	p := &BenchPeer{
		id: BenchIdentifier{id: id},
	}
	return p
}

func (p *BenchPeer) Status() peer.Status {
	return peer.Status{
		PendingRequestCount: int(p.pending.Load()),
		// TODO return real connection status through call back in start/end request
		ConnectionStatus: peer.Available,
	}
}

func (p *BenchPeer) StartRequest() {
	p.pending.Inc()
}

func (p *BenchPeer) EndRequest() {
	p.pending.Dec()
}
