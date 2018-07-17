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
	"go.uber.org/yarpc/peer/pendingheap"
	"go.uber.org/yarpc/peer/roundrobin"
)

type ListFactory func(group *ClientGroup) (peer.ChooserList, error)

var listFactories = make(map[ListType]ListFactory)

func NewRoundRobinList(group *ClientGroup) (peer.ChooserList, error) {
	return roundrobin.New(NewBenchTransport()), nil
}

func NewFewestPendingList(group *ClientGroup) (peer.ChooserList, error) {
	return pendingheap.New(NewBenchTransport()), nil
}

func RegisterList(balancingType ListType, factory ListFactory) error {
	if factory == nil {
		return fmt.Errorf(`unable to register %q, factory is nil`, balancingType)
	}
	if _, ok := listFactories[balancingType]; ok {
		return fmt.Errorf(`factory for %q already exists`, balancingType)
	}
	listFactories[balancingType] = factory
	return nil
}

func InitListFactory() error {
	return multierr.Combine(
		RegisterList(RoundRobin, NewRoundRobinList),
		RegisterList(FewestPending, NewFewestPendingList),
	)
}

func CreateList(group *ClientGroup) (peer.ChooserList, error) {
	factory, ok := listFactories[group.LType]
	if !ok {
		return nil, fmt.Errorf(`type %q is not supported`, group.LType)
	}
	return factory(group)
}
