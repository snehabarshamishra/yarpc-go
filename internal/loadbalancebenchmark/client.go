package loadbalancingbenchmark

import (
	"context"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/peer"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	id        int
	counter   int
	chooser   *peer.BoundChooser
	start     chan EmptySignal
	stop      chan EmptySignal
	wg        *sync.WaitGroup
	sg        *ServerListenerGroup
	sleepTime time.Duration
	lastIssue time.Time
}

func NewClient(id int, group *ClientGroup, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup) (*Client, error) {
	plc, err := CreatePeerListChooser(group, sg.n)
	if err != nil {
		return nil, err
	}
	return &Client{
		id:        id,
		chooser:   plc,
		start:     start,
		stop:      stop,
		wg:        wg,
		sg:        sg,
		sleepTime: time.Second / time.Duration(group.Rps),
	}, nil
}

func (c *Client) issue() (retErr error) {
	res := make(ResponseWriter)
	// concurrent no time out
	ctx := context.Background()
	p, f, err := c.chooser.Choose(ctx, &transport.Request{})
	defer f(retErr)
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(p.Identifier())
	if err != nil {
		return err
	}
	req := c.sg.GetListener(pid)
	req <- res
	<-res
	return err
}

func (c *Client) Start() {
	<-c.start
	//fmt.Println(fmt.Sprintf("client %d start", c.id))
	// use little precision lost to trade off speed
	for {
		select {
		case <-c.stop:
			//fmt.Println(fmt.Sprintf("client %d end", c.id))
			c.wg.Done()
			return
		default:
			c.lastIssue = time.Now()
			go c.issue()
			c.counter += 1
			ct := time.Now()
			time.Sleep(c.sleepTime - time.Duration(ct.Sub(c.lastIssue)))
		}
	}
}
