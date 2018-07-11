package loadbalancingbenchmark

import (
	"fmt"
	"go.uber.org/atomic"
	"time"
)

type Client struct {
	id      int
	counter atomic.Int32
	conn    RequestWriter
	stop    chan EmptySignal
}

func NewClient(id int, conn RequestWriter, stop chan EmptySignal) (*Client, error) {
	return &Client{
		id:   id,
		conn: conn,
		stop: stop,
	}, nil
}

func (c *Client) Generate() {
	for {
		select {
		case <-c.stop:
			fmt.Println(fmt.Sprintf("client %d stopped. ", c.id))
			return
		default:
			res := make(ResponseWriter)
			c.conn <- res
			<-res
			fmt.Println(fmt.Sprintf("client %d received response", c.id))
			time.Sleep(time.Millisecond)
		}
	}
}
