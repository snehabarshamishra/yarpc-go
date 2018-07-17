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

package loadbalancebenchmark

import (
	"fmt"
)

type EmptySignal struct{}

type ResponseWriter chan EmptySignal

type RequestWriter chan ResponseWriter

type ServerListenerGroup struct {
	n         int
	listeners []RequestWriter
}

func NewServerListenerGroup(n int) *ServerListenerGroup {
	var listeners []RequestWriter
	for i := 0; i < n; i++ {
		listeners = append(listeners, make(RequestWriter))
	}
	return &ServerListenerGroup{
		n:         n,
		listeners: listeners,
	}
}

func (sg *ServerListenerGroup) GetListener(pid int) RequestWriter {
	if pid < 0 || pid >= sg.n {
		panic(fmt.Sprintf("pid index out of range, pid: %d size: %d", pid, sg.n))
	}
	return sg.listeners[pid]
}
