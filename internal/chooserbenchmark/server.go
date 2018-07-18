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
	"math/rand"
	"sync"
	"time"
)

type Server struct {
	groupName     *string
	id            int
	counter       int
	listener      RequestWriter
	start         chan struct{}
	stop          chan struct{}
	wg            *sync.WaitGroup
	latencyConfig *LatencyConfig
}

func NewServer(id int, groupName *string, latencyConfig *LatencyConfig, lis RequestWriter, start, stop chan struct{}, wg *sync.WaitGroup) (*Server, error) {
	return &Server{
		groupName:     groupName,
		id:            id,
		listener:      lis,
		start:         start,
		stop:          stop,
		wg:            wg,
		latencyConfig: latencyConfig,
	}, nil
}

func (s *Server) getRandomDelay() time.Duration {
	r := rand.Intn(101)
	if r <= 50 {
		return time.Duration(int(float32(s.latencyConfig.P50) / float32(50) * float32(r)))
	} else if r <= 90 {
		return s.latencyConfig.P50 + time.Duration(int(float32(s.latencyConfig.P90-s.latencyConfig.P50)/float32(40)*float32(r-50)))
	} else if r <= 99 {
		return s.latencyConfig.P90 + time.Duration(int(float32(s.latencyConfig.P99-s.latencyConfig.P90)/float32(9)*float32(r-90)))
	} else {
		return s.latencyConfig.P100
	}
}

func (s *Server) handle(res ResponseWriter) {
	time.Sleep(s.getRandomDelay())
	close(res)
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
