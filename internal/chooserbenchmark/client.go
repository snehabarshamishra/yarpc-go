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

type Client struct {
	groupName   *string
	id          int
	rps         int
	chooser     *peer.BoundChooser
	start       chan struct{}
	stop        chan struct{}
	wg          *sync.WaitGroup
	listeners   *Listeners
	reqCounters []atomic.Int32
	resCounters []atomic.Int32
}

func NewClient(id int, group *ClientGroup, listeners *Listeners, start, stop chan struct{}, wg *sync.WaitGroup, f *PeerListChooserFactory, serverCount int) (*Client, error) {
	plc, err := f.CreatePeerListChooser(group, listeners.n)
	resCounters := make([]atomic.Int32, serverCount)
	reqCounters := make([]atomic.Int32, serverCount)
	if err != nil {
		return nil, err
	}
	return &Client{
		id:          id,
		rps:         group.RPS,
		chooser:     plc,
		start:       start,
		stop:        stop,
		wg:          wg,
		listeners:   listeners,
		reqCounters: reqCounters,
		resCounters: resCounters,
	}, nil
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
	c.reqCounters[pid].Inc()
	req <- res
	response := <-res.channel
	c.resCounters[response.serverId].Inc()
	return err
}

func (c *Client) Start() {
	<-c.start

	sleepTime := time.Second / time.Duration(c.rps)
	for {
		go c.issueRequest()

		timer := time.After(sleepTime)
		select {
		case <-c.stop:
			c.wg.Done()
			return
		case <-timer:
		}
	}
}
