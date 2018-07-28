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
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/net/metrics/bucket"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/peer"
)

var (
	// BucketMs use rpc latency buckets in net metrics
	BucketMs = bucket.NewRPCLatency()
	// BucketLen is the number of rpc latency buckets
	BucketLen = len(BucketMs)
)

// Client issues requests to peers returned by the peer list chooser
type Client struct {
	// identifiers
	groupName string
	id        int

	// metrics
	reqCounter atomic.Int32
	resCounter atomic.Int32
	histogram  []atomic.Int32

	// parameters for normal distribution to increase randomness
	mu    float64
	sigma float64

	// chooser contains peer list implementation and peer list updater
	chooser *peer.BoundChooser

	// each client has a reference for all listeners, index by integer
	listeners Listeners

	start chan struct{}
	stop  chan struct{}
	wg    *sync.WaitGroup
}

// PeerListChooser takes a peer list constructor and server count, returns
// a peer.BoundChooser used by client to pick up peers before issuing requests
func PeerListChooser(constructor PeerListConstructor, serverCount int) *peer.BoundChooser {
	return peer.Bind(constructor(NewBenchTransport()), peer.BindPeers(NewPeerIdentifiers(serverCount)))
}

// NewClient creates a new client
func NewClient(
	id int,
	group *ClientGroup,
	listeners Listeners,
	start, stop chan struct{},
	wg *sync.WaitGroup,
	constructor PeerListConstructor,
	serverCount int,
) *Client {
	plc := PeerListChooser(constructor, serverCount)
	sleepTime := float64(time.Second) / float64(group.RPS)
	return &Client{
		groupName: group.Name,
		id:        id,
		histogram: make([]atomic.Int32, BucketLen),
		mu:        sleepTime,
		sigma:     sleepTime / 20,
		chooser:   plc,
		listeners: listeners,
		start:     start,
		stop:      stop,
		wg:        wg,
	}
}

func (c *Client) normalSleepTime() time.Duration {
	return time.Duration(c.sigma*rand.NormFloat64() + c.mu)
}

func (c *Client) incBucket(t time.Duration) {
	val := int64(t / time.Millisecond)
	if val > 100000 {
		panic(val)
	}
	i := 0
	for i < BucketLen && BucketMs[i] < val {
		i++
	}
	if i == BucketLen {
		i = BucketLen - 1
	}
	c.histogram[i].Inc()
}

func (c *Client) issueRequest() (retErr error) {
	// increase counter each time issue a request
	c.reqCounter.Inc()

	req := Request{
		channel:  make(chan Response),
		clientId: c.id,
	}
	ctx := context.Background() // context no time out
	p, onFinish, err := c.chooser.Choose(ctx, &transport.Request{})
	if err != nil {
		return err
	}
	defer onFinish(retErr)

	// get listener for that peer
	pid, err := strconv.Atoi(p.Identifier())
	if err != nil {
		return err
	}
	lis := c.listeners.Listener(pid)

	start := time.Now()
	// issue the request
	lis <- req
	// wait response
	<-req.channel
	end := time.Now()
	// update latency histogram
	c.incBucket(end.Sub(start))

	// increase counter each time receive a response
	c.resCounter.Inc()

	return nil
}

// Start is the long-run go routine issues requests
func (c *Client) Start() {
	<-c.start

	for {
		go c.issueRequest()

		select {
		case <-c.stop:
			c.wg.Done()
			return
		case <-time.After(c.normalSleepTime()):
		}
	}
}
