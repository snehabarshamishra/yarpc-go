package chooserbenchmark

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/yarpc/api/peer"
	"go.uber.org/yarpc/peer/pendingheap"
	"go.uber.org/yarpc/peer/roundrobin"
)

func TestLaunch(t *testing.T) {
	config := &Config{
		ClientGroups: []ClientGroup{
			{
				Name:  "roundrobin",
				Count: 500,
				RPS:   20,
				Constructor: func(t peer.Transport) peer.ChooserList {
					return roundrobin.New(t)
				},
			},
			{
				Name:  "pendingheap",
				Count: 500,
				RPS:   20,
				Constructor: func(t peer.Transport) peer.ChooserList {
					return pendingheap.New(t)
				},
			},
		},
		ServerGroups: []ServerGroup{
			{
				Name:          "normal",
				Count:         50,
				LatencyConfig: time.Millisecond * 10,
			},
		},
		Duration: 1 * time.Second,
	}
	assert.NoError(t, Run(config))
}
