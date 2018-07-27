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
	"time"
)

type Server struct {
	groupName string
	id        int
	listener  Listener
	start     chan struct{}
	stop      chan struct{}
	wg        *sync.WaitGroup
	latency   *LogNormalLatency
	counter   int
}

func NewServer(
	id int,
	groupName string,
	latency time.Duration,
	lis Listener,
	start, stop chan struct{},
	wg *sync.WaitGroup,
) (*Server, error) {
	return &Server{
		groupName: groupName,
		id:        id,

		listener: lis,
		latency:  NewLogNormalLatency(latency),

		start: start,
		stop:  stop,
		wg:    wg,
	}, nil
}

func (s *Server) handle(res ResponseWriter) {
	time.Sleep(s.latency.Random())
	res.channel <- Message{serverId: s.id}
	close(res.channel)
}

func (s *Server) Serve() {
	<-s.start
	s.wg.Done()
	for {
		select {
		case res := <-s.listener:
			s.counter++
			go s.handle(res)
		case <-s.stop:
			s.wg.Done()
			return
		}
	}
}
