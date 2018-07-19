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

func metric(ctx *Context, ticker *time.Ticker) {
	snapshot := make([][][]int32, 101)
	for i := range snapshot {
		snapshot[i] = make([][]int32, ctx.ClientCount)
		for j := 0; j < ctx.ClientCount; j++ {
			snapshot[i][j] = make([]int32, ctx.ServerCount)
		}
	}
	i := 0
	for {
		select {
		case <-ticker.C:
			for j := 0; j < ctx.ClientCount; j++ {
				for k := 0; k < ctx.ServerCount; k++ {
					snapshot[i][j][k] = ctx.Clients[j].reqCounters[k].Load() - ctx.Clients[j].resCounters[k].Load()
				}
			}
			i++
		case <-ctx.Stop:
			ticker.Stop()
			i = 0
			for i := 0; i < 100; i++ {
				fmt.Println(fmt.Sprintf("show snapshot for %d second", i))
				for j := 0; j < ctx.ClientCount; j++ {
					fmt.Print(fmt.Sprintf("client %d: ", j))
					for k := 0; k < ctx.ServerCount; k++ {
						fmt.Print(fmt.Sprintf("%d ", snapshot[i][j][k]))
					}
				}
				fmt.Println()
			}
			return
		}
	}
}

func launch(ctx *Context) error {
	serverCount := ctx.ServerCount
	clientCount := ctx.ClientCount

	fmt.Println(fmt.Sprintf("launch %d servers...", serverCount))
	for _, server := range ctx.Servers {
		go server.Serve()
	}
	ctx.WG.Add(serverCount)
	close(ctx.ServerStart)
	ctx.WG.Wait()

	fmt.Println(fmt.Sprintf("launch %d clients...", clientCount))
	for _, client := range ctx.Clients {
		go client.Start()
	}
	close(ctx.ClientStart)

	fmt.Println(fmt.Sprintf("launch a backend thread to monitor metrics"))
	ticker := time.NewTicker(ctx.Duration / time.Duration(100))
	go metric(ctx, ticker)

	fmt.Println(fmt.Sprintf("begin benchmark, over after %d seconds...", ctx.Duration/time.Second))
	ctx.WG.Add(serverCount + clientCount)
	time.Sleep(ctx.Duration)
	close(ctx.Stop)
	ctx.WG.Wait()

	return nil
}

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

	time.Sleep(time.Second)
	Visualize(ctx)

	fmt.Println("\nmain workflow is over")
	return nil
}
