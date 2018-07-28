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
)

// Listeners keeps a list of go channels as end points that receive requests
// from clients, it's a shared object among all go routines
type Listeners []Listener

// NewListeners makes n go channels and returns it as Listeners object
func NewListeners(n int) Listeners {
	listeners := make([]Listener, n)
	for i := 0; i < n; i++ {
		listeners[i] = make(Listener)
	}
	return Listeners(listeners)
}

// Listener return the Listener object with corresponding peer id
func (sg Listeners) Listener(pid int) Listener {
	if pid < 0 || pid >= len(sg) {
		panic(fmt.Sprintf("pid index out of range, pid: %d size: %d", pid, len(sg)))
	}
	return sg[pid]
}
