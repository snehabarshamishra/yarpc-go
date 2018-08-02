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
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/peer"
)

const (
	ready = 1
)

// Client issues requests to peers returned by the peer list chooser
type Client struct {
	// identifiers
	groupName string
	id        int

	// metrics
	reqCounter atomic.Int64
	resCounter atomic.Int64
	histogram  *Histogram

	// use normal distribution sleep time to increase randomness
	sleeper *NormalDistSleepTime

	// chooser contains peer list implementation and peer list updater
	chooser *peer.BoundChooser

	// each client has a reference for all listeners, index by integer
	listeners Listeners

	status atomic.Int32

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
) *Client {
	plc := PeerListChooser(group.Constructor, len(listeners))
	return &Client{
		groupName: group.Name,
		id:        id,
		histogram: NewHistogram(BucketMs, int64(time.Millisecond)),
		sleeper:   NewNormalDistSleepTime(group.RPS),
		chooser:   plc,
		listeners: listeners,
		status:    atomic.Int32{},
		start:     start,
		stop:      stop,
		wg:        wg,
	}
}

func (c *Client) normalSleepTime() time.Duration {
	return c.sleeper.Random()
}

func (c *Client) issueRequest() (retErr error) {
	// increase counter each time issue a request
	c.reqCounter.Inc()

	req := Request{
		channel:  make(chan Response),
		clientID: c.id,
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
	if c.status.Load() == ready {
		end := time.Now()
		// update latency histogram
		c.histogram.IncBucket(int64(end.Sub(start)))
		// increase counter each time receive a response
		c.resCounter.Inc()
	}

	return nil
}

// Start is the long-run go routine issues requests
func (c *Client) Start() {
	<-c.start
	c.status.Inc()
	for {
		go func() {
			if err := c.issueRequest(); err != nil {
				panic(err)
			}
		}()

		select {
		case <-c.stop:
			c.wg.Done()
			c.status.Inc()
			return
		case <-time.After(c.normalSleepTime()):
		}
	}
}
