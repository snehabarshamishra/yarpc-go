// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package chooserbenchmark

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckClientGroup(t *testing.T) {
	tests := []struct {
		msg       string
		config    *Config
		wantError string
	}{
		{
			msg: "normal client configuration",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ClientGroups: []ClientGroup{
					{
						Name:  "roundrobin",
						Count: 500,
						RPS:   20,
					},
					{
						Name:  "pendingheap",
						Count: 500,
						RPS:   0,
					},
				},
			},
		},
		{
			msg: "empty client group name",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ClientGroups: []ClientGroup{
					{},
				},
			},
			wantError: "client group name is nil",
		},
		{
			msg: "duplicate client group name",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ClientGroups: []ClientGroup{
					{
						Name:  "foo",
						Count: 1,
					},
					{
						Name:  "foo",
						Count: 1,
					},
				},
			},
			wantError: "client group name duplicated",
		},
		{
			msg: "RPS smaller than 0",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ClientGroups: []ClientGroup{
					{
						Name: "foo",
						RPS:  -1,
					},
				},
			},
			wantError: "RPS field must be greater than 0",
		},
		{
			msg: "group counter smaller than 1",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ClientGroups: []ClientGroup{
					{
						Name:  "foo",
						RPS:   10,
						Count: 0,
					},
				},
			},
			wantError: "number of clients must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			if tt.wantError != "" {
				assert.Contains(t, tt.config.Validate().Error(), tt.wantError)
			} else {
				assert.NoError(t, tt.config.Validate())
			}
		})
	}
}

func TestCheckServerGroup(t *testing.T) {
	tests := []struct {
		msg       string
		config    *Config
		wantError string
	}{
		{
			msg: "normal server configuration",
			config: &Config{
				Duration: time.Second * 10,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{
						Name:          "normal",
						Count:         5,
						LatencyConfig: time.Millisecond * 100,
					},
					{
						Name:          "slow",
						Count:         5,
						LatencyConfig: time.Second,
					},
				},
			},
		},
		{
			msg: "empty server group name",
			config: &Config{
				Duration: time.Second * 10,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{},
				},
			},
			wantError: "server group name is nil",
		},
		{
			msg: "duplicated server group name",
			config: &Config{
				Duration: time.Second * 10,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{
						Name:          "foo",
						Count:         5,
						LatencyConfig: time.Millisecond,
					},
					{
						Name: "foo",
					},
				},
			},
			wantError: "server group name duplicated",
		},
		{
			msg: "latency smaller than 0",
			config: &Config{
				Duration: time.Second * 10,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{
						Name:          "foo",
						Count:         5,
						LatencyConfig: time.Duration(-1),
					},
				},
			},
			wantError: "latency must not be smaller 0",
		},
		{
			msg: "latency greater than duration",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{
						Name:          "foo",
						Count:         5,
						LatencyConfig: 2 * time.Second,
					},
				},
			},
			wantError: "latency must be smaller than test duration",
		},
		{
			msg: "group counter smaller than 1",
			config: &Config{
				Duration: time.Second,
				Output:   os.DevNull,
				ServerGroups: []ServerGroup{
					{
						Name:          "foo",
						Count:         0,
						LatencyConfig: time.Millisecond,
					},
				},
			},
			wantError: "number of servers must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			if tt.wantError != "" {
				assert.Contains(t, tt.config.Validate().Error(), tt.wantError)
			} else {
				assert.NoError(t, tt.config.Validate())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	configZero := &Config{Duration: time.Duration(0)}
	configNegative := &Config{Duration: time.Duration(-1)}
	configNormal := &Config{Duration: time.Duration(100)}
	assert.Contains(t, configZero.Validate().Error(), "test duration must be greater than 0")
	assert.Contains(t, configNegative.Validate().Error(), "test duration must be greater than 0")
	assert.NoError(t, configNormal.Validate())
}
