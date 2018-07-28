// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package chooserbenchmark

import (
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/peer"
)

var _ peer.Peer = (*BenchPeer)(nil)

// BenchPeer is a minimum workable implementation to mock peer behaviors in
// peer-to-peer routing
type BenchPeer struct {
	id      BenchIdentifier
	pending atomic.Int32
	sub     peer.Subscriber
}

// Identifier use BenchIdentifier to generate a unique string
func (p *BenchPeer) Identifier() string {
	return p.id.Identifier()
}

// NewBenchPeer creates a BenchPeer from a integer and peer.Subscriber
func NewBenchPeer(id int, ps peer.Subscriber) *BenchPeer {
	p := &BenchPeer{
		id:  BenchIdentifier{id: id},
		sub: ps,
	}
	return p
}

// Status returns the current status of the BenchPeer.
// pending request count is for policies like fewest pending request
// connection status is always available, no adding peers, peer fails for now
func (p *BenchPeer) Status() peer.Status {
	return peer.Status{
		PendingRequestCount: int(p.pending.Load()),
		ConnectionStatus:    peer.Available,
	}
}

// StartRequest runs at the beginning of a client request
func (p *BenchPeer) StartRequest() {
	p.pending.Inc()
	p.sub.NotifyStatusChanged(p.id)
}

// EndRequest runs in a callback returned by chooser, when request has finished
func (p *BenchPeer) EndRequest() {
	p.pending.Dec()
	p.sub.NotifyStatusChanged(p.id)
}
