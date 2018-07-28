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
	"strconv"

	"go.uber.org/yarpc/api/peer"
)

var _ peer.Identifier = (*BenchIdentifier)(nil)

// BenchIdentifier use integer to uniquely identify a server peer
type BenchIdentifier struct {
	id int
}

// NewPeerIdentifiers create a bunch of peers, id in range [0, n)
func NewPeerIdentifiers(n int) []peer.Identifier {
	if n <= 0 {
		n = 0
	}
	ids := make([]peer.Identifier, n)
	for i := 0; i < n; i++ {
		ids[i] = BenchIdentifier{id: i}
	}
	return ids
}

// Identifier return unique string that identify the peer
func (p BenchIdentifier) Identifier() string {
	return strconv.Itoa(p.id)
}
