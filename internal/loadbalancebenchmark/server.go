package loadbalancingbenchmark

import (
	"fmt"
	"sync"
)

type EmptySignal struct{}

type ResponseWriter chan EmptySignal

type RequestWriter chan ResponseWriter

type Server struct {
	id       int
	counter  int
	listener RequestWriter
	start    chan EmptySignal
	stop     chan EmptySignal
	wg       *sync.WaitGroup
	sg       *ServerListenerGroup
}

func NewServer(id int, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup, config *TestConfig) (*Server, error) {
	return &Server{
		id:       id,
		listener: sg.GetListener(id),
		start:    start,
		stop:     stop,
		wg:       wg,
		sg:       sg,
	}, nil
}

func (s *Server) handle(res ResponseWriter) {
	close(res)
}

func (s *Server) Serve() {
	<-s.start
	fmt.Println(fmt.Sprintf("server %d start serving", s.id))
	for {
		select {
		case res := <-s.listener:
			s.counter += 1
			go s.handle(res)
		case <-s.stop:
			fmt.Println(fmt.Sprintf("server %d stopped. ", s.id))
			s.wg.Done()
			return
		}
	}
}
func (s *Server) GetCount() int {
	return s.counter
}
