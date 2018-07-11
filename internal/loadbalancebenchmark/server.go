package loadbalancingbenchmark

import (
	"fmt"
	"go.uber.org/atomic"
)

type EmptySignal struct{}

type ResponseWriter chan EmptySignal
type RequestWriter chan ResponseWriter

type Server struct {
	id       int
	counter  atomic.Int32
	listener RequestWriter
	stop     chan EmptySignal
}

func NewServer(id int, listener RequestWriter, stop chan EmptySignal) (*Server, error) {
	return &Server{
		id:       id,
		listener: listener,
		stop:     stop,
	}, nil
}

func (s *Server) Serve() {
	for {
		select {
		case res := <-s.listener:
			fmt.Println(fmt.Sprintf("server %d received request", s.id))
			close(res)
		case <-s.stop:
			fmt.Println(fmt.Sprintf("server %d stopped. ", s.id))
			return
		}
	}
}
