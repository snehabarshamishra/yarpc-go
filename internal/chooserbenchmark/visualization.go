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
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tabWidth               = 4
	separator              = "#"
	bar                    = "*"
	displayWidth           = 120
	serverHistogramBuckets = 10
	latencyHeight          = 15
	latencyWidth           = 101
)

type ClientGroupMeta struct {
	name        string
	count       int
	rps         int
	reqCount    int
	resCount    int
	meanLatency time.Duration
	histogram   []int
}

type ServerGroupMeta struct {
	name         string
	servers      []*Server
	histogram    []int
	buckets      []int
	requestCount int
	meanLatency  time.Duration
}

type DisplayConfig struct {
	tabWidth     int
	displayWidth int

	latencyHeight int
	latencyWidth  int

	separator byte
	bar       byte

	serverHistogramBuckets int

	idLen      int
	freqLen    int
	counterLen int

	tableStarUnit     float64
	histogramStarUnit float64
	latencyStarUnit   float64
}

type Visualizer struct {
	clientGroupNames []string
	serverGroupNames []string
	serverData       map[string]*ServerGroupMeta
	clientData       map[string]*ClientGroupMeta

	config *DisplayConfig
}

func NewServerGroupMeta(groupName string, servers []*Server, buckets []int) *ServerGroupMeta {
	requestCount := 0
	counters := make([]int, len(servers))
	for i, server := range servers {
		requestCount += server.counter
		counters[i] = server.counter
	}
	sort.Ints(counters)

	histogram := make([]int, serverHistogramBuckets)
	idx := 0
	for _, counter := range counters {
		for buckets[idx] < counter {
			idx++
		}
		histogram[idx]++
	}

	return &ServerGroupMeta{
		name:         groupName,
		servers:      servers,
		histogram:    histogram,
		buckets:      buckets,
		requestCount: requestCount,
		meanLatency:  servers[0].latency.median,
	}
}

func NewClientGroupMeta(groupName string, clients []*Client) *ClientGroupMeta {
	rps := int(float64(time.Second) / clients[0].mu)
	reqCount, resCount := 0, 0
	histogram := make([]int, BucketLen)
	for _, client := range clients {
		reqCount += int(client.reqCounter.Load())
		resCount += int(client.resCounter.Load())
		for i, freq := range client.histogram {
			histogram[i] += int(freq.Load())
		}
	}
	latencySum := float64(0)
	latencyCount := 0
	for i := 0; i < 50; i++ {
		latencyCount += histogram[i]
		latencySum += float64(histogram[i]) * float64(BucketMs[i])
	}
	return &ClientGroupMeta{
		name:        groupName,
		count:       len(clients),
		rps:         rps,
		reqCount:    reqCount,
		resCount:    resCount,
		meanLatency: time.Duration(latencySum/float64(latencyCount)) * time.Millisecond,
		histogram:   histogram,
	}
}

func NewVisulizer(ctx *Context) *Visualizer {
	servers := ctx.Servers
	serverGroupNames := make([]string, 0)
	serversByGroup := make(map[string][]*Server)
	maxRequestCount, minRequestCount, maxServerCount, maxHistogramVal := 0, math.MaxInt64, 0, 0
	for _, server := range servers {
		if server.counter > maxRequestCount {
			maxRequestCount = server.counter
		}
		if server.counter < minRequestCount {
			minRequestCount = server.counter
		}
		groupName := server.groupName
		if _, ok := serversByGroup[groupName]; !ok {
			serverGroupNames = append(serverGroupNames, groupName)
			serversByGroup[groupName] = make([]*Server, 0)
		}
		serversByGroup[groupName] = append(serversByGroup[groupName], server)
	}
	sort.Strings(serverGroupNames)

	buckets := make([]int, serverHistogramBuckets)
	diff := maxRequestCount - minRequestCount
	width := diff / (serverHistogramBuckets - 1)
	if width == 0 {
		width = 1
	}
	cur := minRequestCount
	for i := 0; i < serverHistogramBuckets; i++ {
		cur += width
		buckets[i] = cur
	}

	serverData := make(map[string]*ServerGroupMeta)
	for _, groupName := range serverGroupNames {
		meta := NewServerGroupMeta(groupName, serversByGroup[groupName], buckets)
		count := len(meta.servers)
		if count > maxServerCount {
			maxServerCount = count
		}
		for _, val := range meta.histogram {
			if val > maxHistogramVal {
				maxHistogramVal = val
			}
		}
		serverData[groupName] = meta
	}

	idLen := NormalizeLength(GetNumLength(maxServerCount))
	counterLen := NormalizeLength(GetNumLength(maxRequestCount))
	freqLen := NormalizeLength(GetNumLength(maxHistogramVal))

	tableStarLen := displayWidth - idLen - counterLen - 2*tabWidth
	tableStarUnit := float64(tableStarLen) / float64(maxRequestCount)

	histogramStarLen := displayWidth - counterLen - freqLen - 2*tabWidth
	histogramStarUnit := float64(histogramStarLen) / float64(maxHistogramVal)

	clients := ctx.Clients
	clientGroupNames := make([]string, 0)
	clientsByGroup := make(map[string][]*Client)
	for _, client := range clients {
		groupName := client.groupName
		if _, ok := clientsByGroup[groupName]; !ok {
			clientGroupNames = append(clientGroupNames, groupName)
			clientsByGroup[groupName] = make([]*Client, 0)
		}
		clientsByGroup[groupName] = append(clientsByGroup[groupName], client)
	}
	sort.Strings(clientGroupNames)

	clientData := make(map[string]*ClientGroupMeta)
	maxLatencyFrequency := 0
	for _, groupName := range clientGroupNames {
		meta := NewClientGroupMeta(groupName, clientsByGroup[groupName])
		for _, freq := range meta.histogram {
			if freq > maxLatencyFrequency {
				maxLatencyFrequency = freq
			}
		}
		clientData[groupName] = meta
	}

	latencyStarUnit := float64(latencyHeight) / float64(maxLatencyFrequency)

	return &Visualizer{
		clientGroupNames: clientGroupNames,
		clientData:       clientData,
		serverGroupNames: serverGroupNames,
		serverData:       serverData,
		config: &DisplayConfig{
			idLen:             idLen,
			freqLen:           freqLen,
			counterLen:        counterLen,
			tableStarUnit:     tableStarUnit,
			histogramStarUnit: histogramStarUnit,
			latencyStarUnit:   latencyStarUnit,
		},
	}
}

func (sgm *ServerGroupMeta) VisualizeServerGroup(vis *Visualizer) {
	config := vis.config
	servers, buckets, histogram := sgm.servers, sgm.buckets, sgm.histogram
	idLen, freqLen, counterLen := config.idLen, config.freqLen, config.counterLen
	tableStarUnit, histogramStarUnit := config.tableStarUnit, config.histogramStarUnit
	name, count, latency, total := sgm.name, len(servers), sgm.meanLatency, sgm.requestCount
	separateLine()

	fmt.Printf(`request count histogram of server group %q`+"\n", name)
	fmt.Printf("number of servers: %d, latency: %v, total requests received: %d\n", count, latency, total)

	if count <= serverHistogramBuckets {
		// if you only have less than 10 servers in this group, just display them one by one
		fmt.Printf("\n%*s\t%*s\t%s\n", idLen, "id", counterLen, "reqs", "histogram")
		for _, server := range servers {
			counter := server.counter
			stars := int(float64(counter) * tableStarUnit)
			fmt.Printf("%*d\t%*d\t%s\n", idLen, server.id, counterLen, server.counter, strings.Repeat(bar, stars))
		}
	} else {
		// when you have more than hundreds of servers, display them in a histogram
		fmt.Printf("\n%*s\t%*s\t%s\n", counterLen, "reqs", freqLen, "freq", "histogram")
		for i, freq := range histogram {
			stars := int(float64(freq) * histogramStarUnit)
			fmt.Printf("%*d\t%*d\t%s\n", counterLen, buckets[i], freqLen, freq, strings.Repeat(bar, stars))
		}
	}
	separateLine()
}

func (cgm *ClientGroupMeta) VisualizeClientGroup(vis *Visualizer) {
	config := vis.config
	name := cgm.name
	histogram := cgm.histogram
	count, rps, reqCount, resCount, meanLatency := cgm.count, cgm.rps, cgm.reqCount, cgm.resCount, cgm.meanLatency
	separateLine()

	fmt.Printf(`request latency histogram of client group %q`+"\n", name)
	fmt.Printf("number of clients: %d, rps: %d, request issued: %d, response received: %d mean latency: %v\n", count, rps, reqCount, resCount, meanLatency)

	fmt.Println()
	pixels := make([][]byte, latencyHeight)
	for i := 0; i < latencyHeight; i++ {
		pixels[i] = make([]byte, latencyWidth)
		pixels[i][0] = '|'
		for j := 1; j < latencyWidth; j++ {
			pixels[i][j] = ' '
		}
	}

	for j := 0; j < BucketLen; j++ {
		stars := int(float64(histogram[j]) * config.latencyStarUnit)
		for i := 0; i < stars; i++ {
			pixels[latencyHeight-1-i][2*(j+1)] = '*'
		}
	}

	for i := 0; i < latencyHeight; i++ {
		for j := 0; j < latencyWidth; j++ {
			fmt.Printf("%c", pixels[i][j])
		}
		fmt.Println()
	}
	fmt.Printf("%s\n", strings.Repeat("-", latencyWidth))

	numLen := 10
	showCount := latencyWidth / numLen
	fmt.Print(" ")
	for i := 1; i <= showCount; i++ {
		fmt.Printf("%*d", numLen, BucketMs[5*i-1])
	}
	fmt.Println()
	separateLine()
}

func Visualize(ctx *Context) {
	fmt.Println("\nbenchmark end, collect metrics and visualize...")

	vis := NewVisulizer(ctx)

	for _, groupName := range vis.clientGroupNames {
		meta := vis.clientData[groupName]
		meta.VisualizeClientGroup(vis)
	}

	for _, groupName := range vis.serverGroupNames {
		meta := vis.serverData[groupName]
		meta.VisualizeServerGroup(vis)
	}
}

func GetNumLength(num int) int {
	return len(strconv.Itoa(num))
}

func NormalizeLength(l int) int {
	return (l + tabWidth - 1) / 4 * 4
}

func separateLine() {
	fmt.Printf("\n%s\n", strings.Repeat(separator, displayWidth))
}
