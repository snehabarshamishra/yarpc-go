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
	"go.uber.org/multierr"
	"sync"
	"time"
)

type Context struct {
	ServerCount int
	Listeners   Listeners
	Servers     []*Server

	ClientCount int
	Clients     []*Client

	WG          sync.WaitGroup
	ServerStart chan struct{}
	ClientStart chan struct{}
	Stop        chan struct{}

	Duration   time.Duration
	MaxLatency time.Duration
}

func (ctx *Context) buildServers(config *Config) error {
	ctx.Servers = make([]*Server, ctx.ServerCount)
	id := 0
	for _, group := range config.ServerGroups {
		for i := 0; i < group.Count; i++ {
			server, err := NewServer(id, group.Name, group.LatencyConfig, ctx.Listeners.Listener(id), ctx.ServerStart, ctx.Stop, &ctx.WG)
			if err != nil {
				return err
			}
			ctx.Servers[id] = server
			if server.latency.median > ctx.MaxLatency {
				ctx.MaxLatency = server.latency.median
			}
			id++
		}
	}
	return nil
}

func (ctx *Context) buildClients(config *Config) error {
	ctx.Clients = make([]*Client, ctx.ClientCount)
	id := 0
	s := time.Now()
	var wg sync.WaitGroup
	total := 0
	for _, group := range config.ClientGroups {
		total += group.Count
	}
	wg.Add(total)
	// it's a CPU bound problem, use 8 go routines can just keep them busy
	for _, group := range config.ClientGroups {
		for i := 0; i < group.Count; i++ {
			client := NewClient(id, &group, ctx.Listeners, ctx.ClientStart, ctx.Stop, &ctx.WG, group.Constructor, ctx.ServerCount)
			ctx.Clients[id] = client
			// Start will append all peers to list, so it's O(ServerCount) complexity.
			go func() {
				client.chooser.Start()
				wg.Done()
			}()
			id++
		}
	}
	wg.Wait()
	e := time.Now()
	fmt.Printf("build %d clients with %d servers in %v\n", total, ctx.ServerCount, e.Sub(s))
	return nil
}

func BuildContext(config *Config) (*Context, error) {
	ctx := Context{
		Duration:    config.Duration,
		ServerStart: make(chan struct{}),
		ClientStart: make(chan struct{}),
		Stop:        make(chan struct{}),
	}

	for _, group := range config.ServerGroups {
		ctx.ServerCount += group.Count
	}
	for _, group := range config.ClientGroups {
		ctx.ClientCount += group.Count
	}

	ctx.Listeners = NewListeners(ctx.ServerCount)

	if err := multierr.Combine(ctx.buildServers(config), ctx.buildClients(config)); err != nil {
		return nil, err
	}

	return &ctx, nil
}
