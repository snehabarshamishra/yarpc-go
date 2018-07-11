package loadbalancingbenchmark

import (
	"fmt"
	"time"
)

func StartBenchmark(config *TestConfig) (retErr error) {
	// register factories
	Init()

	// create n servers
	stop := make(chan EmptySignal)

	lis1 := make(RequestWriter)
	lis2 := make(RequestWriter)
	srv1, _ := NewServer(0, lis1, stop)
	srv2, _ := NewServer(1, lis2, stop)
	go srv1.Serve()
	go srv2.Serve()
	cli, _ := NewClient(0, lis1, stop)
	go cli.Generate()
	time.Sleep(time.Second)
	close(stop)
	time.Sleep(time.Second)
	fmt.Println("main workflow over")
	//chooser, err := CreatePeerChooser(config)
	//if err != nil {
	//	panic(err)
	//}
	//chooser.Start()
	return nil
}
