package loadbalancingbenchmark

import (
	"fmt"
)

type ServerListenerGroup struct {
	n         int
	listeners []RequestWriter
}

func NewServerListenerGroup(n int) *ServerListenerGroup {
	var listeners []RequestWriter
	for i := 0; i < n; i++ {
		listeners = append(listeners, make(RequestWriter))
	}
	return &ServerListenerGroup{
		n:         n,
		listeners: listeners,
	}
}

func (sg *ServerListenerGroup) GetListener(pid int) RequestWriter {
	if pid < 0 || pid >= sg.n {
		panic(fmt.Sprintf("pid index out of range, pid: %d size: %d", pid, sg.n))
	}
	return sg.listeners[pid]
}
