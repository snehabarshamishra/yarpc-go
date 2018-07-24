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
	BucketMs  = bucket.NewRPCLatency()
	BucketLen = len(BucketMs)
)

type Client struct {
	groupName  string
	id         int
	reqCounter int
	resCounter atomic.Int32
	histogram  []atomic.Int32
	mu         float64
	sigma      float64
	chooser    *peer.BoundChooser
	start      chan struct{}
	stop       chan struct{}
	wg         *sync.WaitGroup
	listeners  *Listeners
}

func PeerListChooser(constructor PeerListConstructor, serverCount int) *peer.BoundChooser {
	return peer.Bind(constructor(NewBenchTransport()), peer.BindPeers(NewPeerIdentifiers(serverCount)))
}

func NewClient(id int, group *ClientGroup, listeners *Listeners, start, stop chan struct{}, wg *sync.WaitGroup, constructor PeerListConstructor, serverCount int) *Client {
	plc := PeerListChooser(constructor, serverCount)
	sleepTime := float64(time.Second) / float64(group.RPS)
	return &Client{
		id:        id,
		histogram: make([]atomic.Int32, BucketLen),
		mu:        sleepTime,
		sigma:     sleepTime / 20,
		chooser:   plc,
		start:     start,
		stop:      stop,
		wg:        wg,
		listeners: listeners,
	}
}

func (c *Client) normalSleepTime() time.Duration {
	return time.Duration(c.sigma*rand.NormFloat64() + c.mu)
}

func (c *Client) incBucket(t time.Duration) {
	val := int64(t / time.Millisecond)
	i := 0
	for i < BucketLen && BucketMs[i] < val {
		i++
	}
	c.histogram[i].Inc()
}

func (c *Client) issueRequest() (retErr error) {
	res := ResponseWriter{
		channel:  make(chan Message),
		clientId: c.id,
	}
	// context no time out
	ctx := context.Background()
	p, onFinish, err := c.chooser.Choose(ctx, &transport.Request{})
	if err != nil {
		return err
	}
	defer onFinish(retErr)
	pid, err := strconv.Atoi(p.Identifier())
	if err != nil {
		return err
	}
	req := c.listeners.Listener(pid)
	s := time.Now()
	req <- res
	<-res.channel
	e := time.Now()
	c.resCounter.Inc()
	c.incBucket(e.Sub(s))
	return err
}

func (c *Client) Start() {
	<-c.start

	for {
		go c.issueRequest()
		c.reqCounter++

		timer := time.After(c.normalSleepTime())
		select {
		case <-c.stop:
			c.wg.Done()
			return
		case <-timer:
		}
	}
}
