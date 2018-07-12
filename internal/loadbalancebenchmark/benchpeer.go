package loadbalancingbenchmark

import (
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/peer"
)

type BenchPeer struct {
	id      BenchIdentifier
	pending atomic.Int32
	sub     peer.Subscriber
}

func (p *BenchPeer) Identifier() string {
	return p.id.Identifier()
}

func NewBenchPeer(id int, ps peer.Subscriber) *BenchPeer {
	p := &BenchPeer{
		id:  BenchIdentifier{id: id},
		sub: ps,
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
	p.sub.NotifyStatusChanged(p.id)
}

func (p *BenchPeer) EndRequest() {
	p.pending.Dec()
	p.sub.NotifyStatusChanged(p.id)
}
