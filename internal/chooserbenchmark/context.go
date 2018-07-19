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
	"sync"
	"time"

	"go.uber.org/multierr"
)

type Context struct {
	PLCFactory  *PeerListChooserFactory
	ServerCount int
	Listeners   *Listeners
	Servers     []*Server
	ClientCount int
	Clients     []*Client
	WG          sync.WaitGroup
	ServerStart chan struct{}
	ClientStart chan struct{}
	Stop        chan struct{}
	Duration    time.Duration
}

func (ctx *Context) buildFactories(config *Config) error {
	plcFactory, err := NewPeerListChooserFactory()
	if err != nil {
		return err
	}
	for _, list := range config.CustomListTypes {
		if err := plcFactory.RegisterNewListType(list.ListType, list.ListMethod); err != nil {
			return err
		}
	}
	for _, updater := range config.CustomListUpdaterTypes {
		if err := plcFactory.RegisterNewListUpdaterType(updater.ListUpdaterType, updater.UpdaterMethod); err != nil {
			return err
		}
	}
	ctx.PLCFactory = plcFactory
	return nil
}

func (ctx *Context) buildServers(config *Config) error {
	ctx.Servers = make([]*Server, ctx.ServerCount)
	id := 0
	for _, group := range config.ServerGroups {
		for i := 0; i < group.Count; i++ {
			server, err := NewServer(id, group.Name, group.LatencyConfig, ctx.Listeners.Listener(id), ctx.ServerStart, ctx.Stop, &ctx.WG, ctx.ClientCount)
			if err != nil {
				return err
			}
			ctx.Servers[id] = server
			id++
		}
	}
	return nil
}

func (ctx *Context) buildClients(config *Config) error {
	ctx.Clients = make([]*Client, ctx.ClientCount)
	id := 0
	for _, group := range config.ClientGroups {
		for i := 0; i < group.Count; i++ {
			client, err := NewClient(id, &group, ctx.Listeners, ctx.ClientStart, ctx.Stop, &ctx.WG, ctx.PLCFactory, ctx.ServerCount)
			if err != nil {
				return err
			}
			ctx.Clients[id] = client
			client.chooser.Start()
			id++
		}
	}
	return nil
}

func BuildContext(config *Config) (*Context, error) {
	fmt.Println("build test context....")

	ctx := Context{
		Duration:    config.Duration,
		ServerStart: make(chan struct{}),
		ClientStart: make(chan struct{}),
		Stop:        make(chan struct{}),
	}

	if err := ctx.buildFactories(config); err != nil {
		return nil, err
	}

	for _, group := range config.ServerGroups {
		ctx.ServerCount += group.Count
	}
	for _, group := range config.ClientGroups {
		ctx.ClientCount += group.Count
	}

	ctx.Listeners = NewListeners(ctx.ServerCount)

	err := multierr.Combine(
		ctx.buildServers(config),
		ctx.buildClients(config),
	)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}
