package loadbalancingbenchmark

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/peer"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	id      int
	counter atomic.Int32
	chooser *peer.BoundChooser
	start   chan EmptySignal
	stop    chan EmptySignal
	wg      *sync.WaitGroup
	sg      *ServerListenerGroup
}

func NewClient(id int, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup, config *TestConfig) (*Client, error) {
	chooser, err := CreatePeerChooser(config)
	chooser.Start()
	if err != nil {
		return nil, err
	}
	return &Client{
		id:      id,
		chooser: chooser,
		sg:      sg,
		start:   start,
		stop:    stop,
		wg:      wg,
	}, nil
}

func (c *Client) issue() (retErr error) {
	res := make(ResponseWriter)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, f, err := c.chooser.Choose(ctx, &transport.Request{})
	defer f(retErr)
	pid, err := strconv.Atoi(p.Identifier())
	req := c.sg.GetListener(pid)
	req <- res
	<-res
	return err
}

func (c *Client) Start() {
	<-c.start
	fmt.Println(fmt.Sprintf("client %d start", c.id))
	for {
		select {
		case <-c.stop:
			fmt.Println(fmt.Sprintf("client %d stopped. ", c.id))
			c.wg.Done()
			return
		default:
			go c.issue()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
