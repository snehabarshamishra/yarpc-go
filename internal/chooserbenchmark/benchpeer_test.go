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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/yarpc/api/peer"
)

type FakePeerSubscriber struct {
	peers    map[string]*BenchPeer
	counters map[string]int32
}

func NewFakePeerSubscriber() *FakePeerSubscriber {
	return &FakePeerSubscriber{
		peers:    make(map[string]*BenchPeer),
		counters: make(map[string]int32),
	}
}

func (sub *FakePeerSubscriber) Register(pid peer.Identifier, p *BenchPeer) {
	sub.peers[pid.Identifier()] = p
}

func (sub *FakePeerSubscriber) NotifyStatusChanged(pid peer.Identifier) {
	id := pid.Identifier()
	if _, ok := sub.peers[id]; !ok {
		panic(fmt.Sprintf(`peer %q not registered`, id))
	}
	sub.counters[id] = sub.peers[id].pending.Load()
}

func TestNewBenchPeer(t *testing.T) {
	p := NewBenchPeer(0, &FakePeerSubscriber{})
	assert.Equal(t, "0", p.Identifier())
}

func TestStartEndRequest(t *testing.T) {
	sub := NewFakePeerSubscriber()
	p1 := NewBenchPeer(0, sub)
	p2 := NewBenchPeer(1, sub)
	sub.Register(p1.id, p1)
	sub.Register(p2.id, p2)

	p1.StartRequest()
	p2.StartRequest()
	p2.StartRequest()

	assert.Equal(t, int32(1), p1.pending.Load())
	assert.Equal(t, int32(1), sub.counters[p1.id.Identifier()])

	p1.EndRequest()

	assert.Equal(t, int32(0), p1.pending.Load())
	assert.Equal(t, int32(0), sub.counters[p1.id.Identifier()])
	assert.Equal(t, int32(2), p2.pending.Load())
	assert.Equal(t, int32(2), sub.counters[p2.id.Identifier()])
}
