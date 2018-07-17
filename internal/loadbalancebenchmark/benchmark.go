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
	"math/rand"
	"sync"
	"time"
)

func StartBenchmark(config *TestConfig) error {
	fmt.Println("checking config...")
	if err := config.Validate(); err != nil {
		return err
	}

	fmt.Println("start constructing servers and clients...")
	rand.Seed(time.Now().UnixNano())
	InitFactory()

	totalServerCount := config.GetServerTotalCount()
	totalClientCount := config.GetClientTotalCount()
	sg := NewServerListenerGroup(totalServerCount)

	serverStart := make(chan EmptySignal)
	clientStart := make(chan EmptySignal)
	stop := make(chan EmptySignal)

	var wg sync.WaitGroup

	servers, err := createServers(totalServerCount, config, sg, serverStart, stop, &wg)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("created %d servers", totalServerCount))
	clients, err := createClients(totalClientCount, config, sg, clientStart, stop, &wg)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("created %d clients", totalClientCount))

	// launch clients and servers
	fmt.Println("launch servers...")
	for _, server := range servers {
		go server.Serve()
	}
	wg.Add(totalServerCount)
	close(serverStart)
	wg.Wait()

	fmt.Println("launch clients...")
	for _, client := range clients {
		go client.Start()
	}
	close(clientStart)

	fmt.Println(fmt.Sprintf("start benchmark, over after %d seconds...", config.Duration/time.Second))
	wg.Add(totalServerCount + totalClientCount)
	time.Sleep(config.Duration)
	close(stop)
	wg.Wait()

	fmt.Println("benchmark is over, collect metrics and visualize...")
	visualize(totalServerCount, servers, clients)

	fmt.Println("main workflow is over")
	return nil
}

func createServers(total int, config *TestConfig, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup) ([]*Server, error) {
	groups := config.ServerGroup
	servers := make([]*Server, total)
	id := 0
	for _, group := range groups {
		for i := 0; i < group.Cnt; i++ {
			server, err := NewServer(id, group.Type, group.LatencyConfig, sg, start, stop, wg)
			if err != nil {
				return nil, err
			}
			servers[id] = server
			id++
		}
	}
	return servers, nil
}

func createClients(total int, config *TestConfig, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup) ([]*Client, error) {
	groups := config.ClientGroup
	clients := make([]*Client, total)
	id := 0
	for _, group := range groups {
		for i := 0; i < group.Cnt; i++ {
			client, err := NewClient(id, &group, sg, start, stop, wg)
			if err != nil {
				return nil, err
			}
			client.chooser.Start()
			clients[id] = client
			id++
		}
	}
	return clients, nil
}

func visualize(totalServerCount int, servers []*Server, clients []*Client) {
	resTotal := 0
	ids := make([]int, totalServerCount)
	cnts := make([]int, totalServerCount)
	types := make([]MachineType, totalServerCount)
	maxCnt := 0
	for i, server := range servers {
		id := server.id
		ids[i] = id
		c := server.counter
		cnts[i] = c
		if c > maxCnt {
			maxCnt = c
		}
		types[i] = server.machineType
		resTotal += c
		fmt.Println(fmt.Sprintf("server %d, counter: %d", id, c))
	}
	reqTotal := 0
	for _, client := range clients {
		reqTotal += client.counter
	}
	fmt.Println(fmt.Sprintf("server counter: %d, client counter: %d", resTotal, reqTotal))

	PrintToTerminal(ids, cnts, types, maxCnt, totalServerCount)
}
