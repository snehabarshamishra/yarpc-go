package loadbalancingbenchmark

import (
	"math/rand"
	"sync"
	"time"
)

type Server struct {
	id            int
	counter       int
	listener      RequestWriter
	start         chan EmptySignal
	stop          chan EmptySignal
	wg            *sync.WaitGroup
	latencyConfig *LatencyConfig
	machineType   MachineType
}

func NewServer(id int, machineType MachineType, latencyConfig *LatencyConfig, sg *ServerListenerGroup, start, stop chan EmptySignal, wg *sync.WaitGroup) (*Server, error) {
	return &Server{
		id:            id,
		listener:      sg.GetListener(id),
		start:         start,
		stop:          stop,
		wg:            wg,
		latencyConfig: latencyConfig,
		machineType:   machineType,
	}, nil
}

func (s *Server) getRandomDelay() time.Duration {
	r := rand.Intn(101)
	if r <= 50 {
		return time.Duration(int(float32(s.latencyConfig.P50) / float32(50) * float32(r)))
	} else if r <= 90 {
		return s.latencyConfig.P50 + time.Duration(int(float32(s.latencyConfig.P90-s.latencyConfig.P50)/float32(40)*float32(r-50)))
	} else if r <= 99 {
		return s.latencyConfig.P90 + time.Duration(int(float32(s.latencyConfig.P99-s.latencyConfig.P90)/float32(9)*float32(r-90)))
	} else {
		return s.latencyConfig.P100
	}
}

func (s *Server) handle(res ResponseWriter) {
	delay := s.getRandomDelay()
	time.Sleep(delay)
	close(res)
}

func (s *Server) Serve() {
	<-s.start
	s.wg.Done()
	for {
		select {
		case res := <-s.listener:
			s.counter++
			go s.handle(res)
		case <-s.stop:
			s.wg.Done()
			return
		}
	}
}
