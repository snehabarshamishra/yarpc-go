// Copyright (c) 2017 Uber Technologies, Inc.
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

package grpc_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "go.uber.org/yarpc/x/yarpctest"
)

func TestStreaming(t *testing.T) {
	tests := []struct {
		name     string
		services Lifecycle
		requests Action
	}{
		{
			name: "stream requests",
			services: Lifecycles(
				GRPCService(
					Name("myservice"),
					Port(31112),
					Proc(
						Name("proc"),
						EchoStreamHandler(),
					),
				),
			),
			requests: Actions(
				GRPCStreamRequest(
					Port(31112),
					Service("myservice"),
					Procedure("proc"),
					ClientStreamActions(
						SendStreamMsg("test"),
						RecvStreamMsg("test"),
						SendStreamMsg("test2"),
						RecvStreamMsg("test2"),
						CloseStream(),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.services.Start(t))
			defer func() { require.NoError(t, tt.services.Stop(t)) }()
			tt.requests.Run(t)
		})
	}
}
