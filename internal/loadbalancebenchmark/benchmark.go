package loadbalancingbenchmark

import (
	"fmt"
	"sync"
	"time"
)

func StartBenchmark(config *TestConfig) error {
	// register factories
	Init()

	// get parameters from config
	serverCount, err := config.GetServerCount()
	if err != nil {
		return err
	}
	clientCount, err := config.GetClientCount()
	if err != nil {
		return err
	}
	sg := NewServerListenerGroup(serverCount)
	start := make(chan EmptySignal)
	stop := make(chan EmptySignal)

	var wg sync.WaitGroup
	wg.Add(serverCount + clientCount)

	servers, err := createServers(serverCount, sg, start, stop, &wg, config)
	if err != nil {
		return err
	}
	_, err = createClients(clientCount, sg, start, stop, &wg, config)
	if err != nil {
		return err
	}

	close(start)
	time.Sleep(time.Second)
	close(stop)
	wg.Wait()
	fmt.Println("collect metrics")
	for i := 0; i < serverCount; i++ {
		fmt.Println(fmt.Sprintf("server %d request count: %d", i, servers[i].GetCount()))
	}
	fmt.Println("main workflow over")
	return nil
}

func createServers(m int, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup, config *TestConfig) ([]*Server, error) {
	var servers []*Server
	for i := 0; i < m; i++ {
		server, err := NewServer(i, sg, start, stop, wg, config)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
		go server.Serve()
	}
	return servers, nil
}

func createClients(n int, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup, config *TestConfig) ([]*Client, error) {
	var clients []*Client
	for i := 0; i < n; i++ {
		client, err := NewClient(i, sg, start, stop, wg, config)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
		go client.Start()
	}
	return clients, nil
}
