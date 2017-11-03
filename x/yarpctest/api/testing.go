// Copyright (c) 2017 Uber Technologies, Inc.
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

package api

import (
	"sync"
	"testing"
)

// TestingT is an interface wrapper around *testing.T and *testing.B
type TestingT interface {
	testing.TB
}

// Run will cast the TestingT to it's sub and call the appropriate Run func.
func Run(name string, t TestingT, f func(TestingT)) {
	if tt, ok := t.(*testing.T); ok {
		tt.Run(name, func(ttt *testing.T) { f(ttt) })
		return
	}
	if tb, ok := t.(*testing.B); ok {
		tb.Run(name, func(ttb *testing.B) { f(ttb) })
		return
	}
	t.Error("invalid test harness")
	t.FailNow()
}

// TestingTInjectable is an embeddable struct that can pass on TestingT information for
// lifecycles.
type TestingTInjectable struct {
	sync.Mutex
	t TestingT
}

// Start safely sets the TestingT for the testable.
func (b *TestingTInjectable) Start(t TestingT) error {
	b.Lock()
	b.t = t
	b.Unlock()
	return nil
}

// GetTestingT safely gets the TestingT for the testable.
func (b *TestingTInjectable) GetTestingT() TestingT {
	b.Lock()
	t := b.t
	b.Unlock()
	return t
}

// NoopLifecycle can be embedded in a struct to make it implement the Start
// and Stop methods.
type NoopLifecycle struct{}

// Start is a Noop.
func (b *NoopLifecycle) Start(t TestingT) error {
	return nil
}

// Stop is a Noop.
func (b *NoopLifecycle) Stop(t TestingT) error {
	return nil
}

// NoopStop can be embedded in a struct to make it implement the "Stop" method.
type NoopStop struct{}

// Stop is a Noop.
func (b *NoopStop) Stop(t TestingT) error {
	return nil
}
