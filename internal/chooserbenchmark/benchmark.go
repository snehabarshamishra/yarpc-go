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
	"math/rand"
	"time"
)

func launch(ctx *Context) error {
	serverCount := ctx.ServerCount
	clientCount := ctx.ClientCount

	fmt.Printf("launch %d servers...\n", serverCount)
	for _, server := range ctx.Servers {
		go server.Serve()
	}
	// wait until all servers start, ensure all servers are ready when clients
	// begin to issue requests
	ctx.WG.Add(serverCount)
	close(ctx.ServerStart)
	ctx.WG.Wait()

	fmt.Printf("launch %d clients...\n", clientCount)
	for _, client := range ctx.Clients {
		go client.Start()
	}

	fmt.Printf("begin benchmark, over after %d seconds...\n", ctx.Duration/time.Second)
	close(ctx.ClientStart)
	time.Sleep(ctx.Duration)

	// wait until all servers and clients stop
	ctx.WG.Add(serverCount + clientCount)
	close(ctx.Stop)
	ctx.WG.Wait()

	return nil
}

// Run start a benchmark according to your configuration, is the main entry
func Run(config *Config) error {
	// use the seed function to initialize the default source, default source is
	// safe for concurrent use by multiple go routines
	rand.Seed(time.Now().UnixNano())

	if err := config.Validate(); err != nil {
		return err
	}

	ctx, err := BuildContext(config)
	if err != nil {
		return err
	}

	if err := launch(ctx); err != nil {
		return err
	}

	Visualize(ctx)

	return nil
}
