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
	fmt.Println("\nbenchmark end, collect metrics and visualize...")

	reqTotal, rcvTotal := 0, 0
	servers, clients := ctx.Servers, ctx.Clients
	maxCounter := 0

	for _, server := range servers {
		c := 0
		for _, cc := range server.reqCounters {
			c += cc
		}
		if c > maxCounter {
			maxCounter = c
		}
		rcvTotal += c
		fmt.Println(fmt.Sprintf("server %d, counter: %d", server.id, c))
	}

	for _, client := range clients {
		c := int32(0)
		for _, cc := range client.resCounters {
			c += cc.Load()
		}
		reqTotal += int(c)
	}
	fmt.Println(fmt.Sprintf("total request issued: %d, total request received: %d", reqTotal, rcvTotal))

	maxStarCount := 60
	base := float64(maxCounter) / float64(maxStarCount)

	fmt.Println("\nname\t\tcount\t\thistogram")
	for _, server := range servers {
		c := 0
		for _, cc := range server.reqCounters {
			c += cc
		}
		starCount := int(float64(c) / base)
		fmt.Println(fmt.Sprintf("%s\t\t%d\t\t%s", *server.groupName, c, strings.Repeat("*", starCount)))
	}

	fmt.Println("per client received response")
	for _, client := range clients {
		fmt.Println(fmt.Sprintf("client %d response report:", client.id))
		for i := 0; i < ctx.ServerCount; i++ {
			fmt.Println(fmt.Sprintf("received %d responses from server %d", client.resCounters[i].Load(), i))
		}
	}

	fmt.Println("per server received request")
	for _, server := range servers {
		fmt.Println(fmt.Sprintf("server %d request report:", server.id))
		for i := 0; i < ctx.ClientCount; i++ {
			fmt.Println(fmt.Sprintf("received %d requests from client %d", server.reqCounters[i], i))
		}
	}
}
