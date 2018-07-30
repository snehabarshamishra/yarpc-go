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
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// display area size
	tabWidth      = 4
	displayWidth  = 120
	latencyHeight = 15
	latencyWidth  = 101

	// ascii code use to display histogram
	separator = "#"
	bar       = "*"

	// server request counter bucket count
	serverBucketCount = 10
)

// clientGroupMeta contains information to visualize a client group
type clientGroupMeta struct {
	// raw
	name string

	// counter related metrics
	count     int
	reqCount  int
	resCount  int
	histogram []int

	// time related metrics
	rps         int
	meanLatency time.Duration
}

// serverGroupMeta contains information to visualize a server group
type serverGroupMeta struct {
	// raw
	name    string
	servers []*Server

	// counter related metrics
	requestCount int
	buckets      []int
	histogram    []int

	// time related metrics
	meanLatency time.Duration
}

// gauge contains settings to control display in terminal
type gauge struct {
	// field size
	idLen      int
	freqLen    int
	counterLen int

	// indicates one star stands for how many requests or servers
	tableStarUnit     float64
	histogramStarUnit float64
	latencyStarUnit   float64
}

// Visualizer is the visualization module for benchmark
type Visualizer struct {
	clientGroupNames []string
	serverGroupNames []string
	serverData       map[string]*serverGroupMeta
	clientData       map[string]*clientGroupMeta

	gauge *gauge
}

func newServerGroupMeta(groupName string, servers []*Server, buckets []int) *serverGroupMeta {
	requestCount := 0
	counters := make([]int, len(servers))
	for i, server := range servers {
		requestCount += server.counter
		counters[i] = server.counter
	}
	sort.Ints(counters)

	histogram := make([]int, serverBucketCount)
	idx := 0
	for _, counter := range counters {
		for buckets[idx] < counter {
			idx++
		}
		histogram[idx]++
	}

	// all servers have same latency configuration, median of log normal
	// distribution is the mean latency by definition
	meanLatency := servers[0].latency.median

	return &serverGroupMeta{
		name:         groupName,
		servers:      servers,
		requestCount: requestCount,
		buckets:      buckets,
		histogram:    histogram,
		meanLatency:  meanLatency,
	}
}

func newClientGroupMeta(groupName string, clients []*Client) *clientGroupMeta {
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
	return &clientGroupMeta{
		name:        groupName,
		count:       len(clients),
		reqCount:    reqCount,
		resCount:    resCount,
		histogram:   histogram,
		rps:         rps,
		meanLatency: time.Duration(latencySum/float64(latencyCount)) * time.Millisecond,
	}
}

// aggregate servers into a hash table, also returns maximum request count and
// minimum request at the same time
func aggregatedServersByGroupName(servers []*Server) ([]string, map[string][]*Server, int, int) {
	serverGroupNames := make([]string, 0)
	serversByGroup := make(map[string][]*Server)
	maxRequestCount, minRequestCount := 0, 0

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

	return serverGroupNames, serversByGroup, maxRequestCount, minRequestCount
}

// create global request counter bucket based on maximum request count and
// minimum request count
func serverRequestCounterBucket(maxRequestCount, minRequestCount int) []int {
	buckets := make([]int, serverBucketCount)
	diff := maxRequestCount - minRequestCount
	width := diff / (serverBucketCount - 1)
	if width == 0 {
		width = 1
	}
	cur := minRequestCount
	for i := 0; i < serverBucketCount; i++ {
		cur += width
		buckets[i] = cur
	}
	return buckets
}

func populateServerData(
	serverGroupNames []string,
	serversByGroup map[string][]*Server,
	buckets []int,
) (map[string]*serverGroupMeta, int, int) {
	maxServerCount, maxServerHistogramVal := 0, 0
	serverData := make(map[string]*serverGroupMeta)

	for _, groupName := range serverGroupNames {
		meta := newServerGroupMeta(groupName, serversByGroup[groupName], buckets)
		count := len(meta.servers)
		if count > maxServerCount {
			maxServerCount = count
		}
		for _, val := range meta.histogram {
			if val > maxServerHistogramVal {
				maxServerHistogramVal = val
			}
		}
		serverData[groupName] = meta
	}

	return serverData, maxServerCount, maxServerHistogramVal
}

// gaugeCalculation calculated field width and scalar for stars in histograms
// latencyStarUnit will be calculated after getting clientData
func gaugeCalculation(maxServerCount, maxRequestCount, maxServerHistogramVal int) *gauge {
	idLen := normalizeLength(getNumLength(maxServerCount))
	counterLen := normalizeLength(getNumLength(maxRequestCount))
	freqLen := normalizeLength(getNumLength(maxServerHistogramVal))

	tableStarLen := displayWidth - idLen - counterLen - 2*tabWidth
	tableStarUnit := float64(tableStarLen) / float64(maxRequestCount)

	histogramStarLen := displayWidth - counterLen - freqLen - 2*tabWidth
	histogramStarUnit := float64(histogramStarLen) / float64(maxServerHistogramVal)

	return &gauge{
		idLen:             idLen,
		freqLen:           freqLen,
		counterLen:        counterLen,
		tableStarUnit:     tableStarUnit,
		histogramStarUnit: histogramStarUnit,
	}
}

func populateClientData(clients []*Client) ([]string, map[string]*clientGroupMeta, int) {
	clientGroupNames := make([]string, 0)
	clientsByGroup := make(map[string][]*Client)
	clientData := make(map[string]*clientGroupMeta)

	for _, client := range clients {
		groupName := client.groupName
		if _, ok := clientsByGroup[groupName]; !ok {
			clientGroupNames = append(clientGroupNames, groupName)
			clientsByGroup[groupName] = make([]*Client, 0)
		}
		clientsByGroup[groupName] = append(clientsByGroup[groupName], client)
	}
	sort.Strings(clientGroupNames)

	maxLatencyFrequency := 0
	for _, groupName := range clientGroupNames {
		meta := newClientGroupMeta(groupName, clientsByGroup[groupName])
		for _, freq := range meta.histogram {
			if freq > maxLatencyFrequency {
				maxLatencyFrequency = freq
			}
		}
		clientData[groupName] = meta
	}

	return clientGroupNames, clientData, maxLatencyFrequency
}

// NewVisualizer returns a Visualizer for metrics data visualization
func NewVisualizer(ctx *Context) *Visualizer {
	// aggregate servers, like group by GroupName in sql
	serverGroupNames, serversByGroup, maxRequestCount, minRequestCount := aggregatedServersByGroupName(ctx.Servers)
	// calculate request counter buckets based on range
	buckets := serverRequestCounterBucket(maxRequestCount, minRequestCount)
	// populate necessary data for server group visualization
	serverData, maxServerCount, maxServerHistogramVal := populateServerData(serverGroupNames, serversByGroup, buckets)
	// calculate base unit for a star character in terminal
	gauge := gaugeCalculation(maxServerCount, maxRequestCount, maxServerHistogramVal)
	// populate necessary data for client group visualization
	clientGroupNames, clientData, maxLatencyFrequency := populateClientData(ctx.Clients)
	// set base unit for latency graph when get maximum latency frequency
	gauge.latencyStarUnit = float64(latencyHeight) / float64(maxLatencyFrequency)

	return &Visualizer{
		clientGroupNames: clientGroupNames,
		clientData:       clientData,
		serverGroupNames: serverGroupNames,
		serverData:       serverData,
		gauge:            gauge,
	}
}

// visualize a server group
func (sgm *serverGroupMeta) visualizeServerGroup(vis *Visualizer) {
	gauge := vis.gauge
	servers, buckets, histogram := sgm.servers, sgm.buckets, sgm.histogram
	idLen, freqLen, counterLen := gauge.idLen, gauge.freqLen, gauge.counterLen
	tableStarUnit, histogramStarUnit := gauge.tableStarUnit, gauge.histogramStarUnit
	name, count, latency, total := sgm.name, len(servers), sgm.meanLatency, sgm.requestCount
	separateLine()

	fmt.Printf(`request count histogram of server group %q`+"\n", name)
	fmt.Printf("number of servers: %d, latency: %v, total requests received: %d\n", count, latency, total)

	if count <= serverBucketCount {
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

// visualize a client group
func (cgm *clientGroupMeta) visualizeClientGroup(vis *Visualizer) {
	gauge := vis.gauge
	name := cgm.name
	histogram := cgm.histogram
	count, rps, reqCount, resCount, meanLatency := cgm.count, cgm.rps, cgm.reqCount, cgm.resCount, cgm.meanLatency
	separateLine()

	fmt.Printf(`request latency histogram of client group %q`+"\n", name)
	fmt.Printf("number of clients: %d, rps: %d, request issued: %d, response received: %d mean latency: %v\n",
		count, rps, reqCount, resCount, meanLatency)

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
		stars := int(float64(histogram[j]) * gauge.latencyStarUnit)
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

// Visualize do the visualization of metrics data stored in context
func Visualize(ctx *Context) {
	fmt.Println("\nbenchmark end, collect metrics and visualize...")

	vis := NewVisualizer(ctx)

	for _, groupName := range vis.clientGroupNames {
		meta := vis.clientData[groupName]
		meta.visualizeClientGroup(vis)
	}

	for _, groupName := range vis.serverGroupNames {
		meta := vis.serverData[groupName]
		meta.visualizeServerGroup(vis)
	}
}

func getNumLength(num int) int {
	return len(strconv.Itoa(num))
}

func normalizeLength(l int) int {
	return (l + tabWidth - 1) / 4 * 4
}

func separateLine() {
	fmt.Printf("\n%s\n", strings.Repeat(separator, displayWidth))
}
