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
	"strings"
)

func Visualize(ctx *Context) {
	fmt.Println("benchmark end, collect metrics and visualize...")

	reqTotal, rcvTotal := 0, 0
	servers, clients := ctx.Servers, ctx.Clients
	maxCounter := 0

	for _, server := range servers {
		c := server.counter
		if c > maxCounter {
			maxCounter = c
		}
		rcvTotal += c
		fmt.Println(fmt.Sprintf("server %d, counter: %d", server.id, c))
	}

	for _, client := range clients {
		reqTotal += client.counter
	}
	fmt.Println(fmt.Sprintf("total request issued: %d, total request received: %d", reqTotal, rcvTotal))

	maxStarCount := 60
	base := float64(maxCounter) / float64(maxStarCount)

	for _, server := range servers {
		starCount := int(float64(server.counter) / base)
		fmt.Println(fmt.Sprintf("%s\t\t%d\t\t%s", *server.groupName, server.counter, strings.Repeat("*", starCount)))
	}
}
