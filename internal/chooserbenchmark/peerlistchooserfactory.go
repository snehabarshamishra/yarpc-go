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
	"go.uber.org/yarpc/peer"
)

type PeerListChooserFactory struct {
	listFactory        *ListFactory
	listUpdaterFactory *ListUpdaterFactory
}

func NewPeerListChooserFactory() (*PeerListChooserFactory, error) {
	listFactory, err := NewListFacotry()
	if err != nil {
		return nil, err
	}
	listUpdaterFactory, err := NewListUpdaterFactory()
	if err != nil {
		return nil, err
	}
	return &PeerListChooserFactory{
		listFactory:        listFactory,
		listUpdaterFactory: listUpdaterFactory,
	}, nil
}

func (f *PeerListChooserFactory) RegisterNewListType(listType ListType, listMethod ListFactoryMethod) error {
	return f.listFactory.Register(listType, listMethod)
}

func (f *PeerListChooserFactory) RegisterNewListUpdaterType(listUpdaterType ListUpdaterType, listUpdaterMethod ListUpdaterFactoryMethod) error {
	return f.listUpdaterFactory.Register(listUpdaterType, listUpdaterMethod)
}

func (f *PeerListChooserFactory) CreatePeerListChooser(group *ClientGroup, serverCount int) (*peer.BoundChooser, error) {
	list, err := f.listFactory.CreateList(group)
	if err != nil {
		return nil, err
	}
	listUpdater, err := f.listUpdaterFactory.CreateListUpdater(group, serverCount)
	if err != nil {
		return nil, err
	}
	return peer.Bind(list, listUpdater), nil
}
