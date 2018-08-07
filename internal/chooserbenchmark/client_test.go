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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/peer"
	"go.uber.org/yarpc/peer/roundrobin"
)

func TestClient(t *testing.T) {
	// initiate parameters for NewClient
	clientGroup := &ClientGroup{
		Name:        "roundrobin",
		Count:       1,
		RPS:         1000,
		Constructor: func(t peer.Transport) peer.ChooserList { return roundrobin.New(t) },
	}
	listeners := NewListeners(1)
	start, stop := make(chan struct{}), make(chan struct{})
	wg := sync.WaitGroup{}
	// create a new client and start the peer list chooser
	client := NewClient(0, clientGroup, listeners, start, stop, &wg)
	client.chooser.Start()
	lis, err := listeners.Listener(0)
	assert.NoError(t, err)

	reqCounter := atomic.Int64{}
	// start client go routine
	go client.Start()
	// start server go routine
	go func(lis Listener) {
		for {
			select {
			case req := <-lis:
				close(req.channel)
				reqCounter.Inc()
			case <-stop:
				wg.Done()
				return
			}
		}
	}(lis)
	assert.Equal(t, int64(0), reqCounter.Load(), "shouldn't receive request before client start")

	// start client and server
	close(start)
	time.Sleep(time.Millisecond * 10)

	// stop client and server
	wg.Add(2)
	close(stop)
	wg.Wait()

	resCount1 := client.resCounter.Load()
	assert.True(t, reqCounter.Load() > 0 && resCount1 <= reqCounter.Load(),
		"request received by server should greater than or equal to response received by clients")
	// sleep another 10 milliseconds to test whether we get any response after test is over
	time.Sleep(time.Millisecond * 10)
	resCount2 := client.resCounter.Load()
	assert.True(t, resCount1 == resCount2, "shouldn't receive any response after test is over")
	// close listener
	close(lis)
}
