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

	"go.uber.org/multierr"
	"go.uber.org/yarpc/api/peer"
	peer2 "go.uber.org/yarpc/peer"
)

type ListUpdaterFactory func(serverCount int) (peer.Binder, error)

var listUpdaterFactories = make(map[ListUpdaterType]ListUpdaterFactory)

func NewStaticListUpdater(serverCount int) (peer.Binder, error) {
	return peer2.BindPeers(MakePeerIdentifier(serverCount)), nil
}

func RegisterListUpdater(updaterType ListUpdaterType, factory ListUpdaterFactory) error {
	if factory == nil {
		return fmt.Errorf(`unable to register %q, factory is nil`, updaterType)
	}
	if _, ok := listUpdaterFactories[updaterType]; ok {
		return fmt.Errorf(`factory for %q already exists`, updaterType)
	}
	listUpdaterFactories[updaterType] = factory
	return nil
}

func InitListUpdaterFactory() error {
	return multierr.Combine(
		RegisterListUpdater(Static, NewStaticListUpdater),
	)
}

func CreateListUpdater(group *ClientGroup, serverCount int) (peer.Binder, error) {
	factory, ok := listUpdaterFactories[group.LUType]
	if !ok {
		return nil, fmt.Errorf(`type %q is not supported`, group.LUType)
	}
	return factory(serverCount)
}
