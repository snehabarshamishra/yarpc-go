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

package yarpctest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/yarpc/yarpcerrors"
)

func TestServiceRouting(t *testing.T) {
	tests := []struct {
		name     string
		services Lifecycle
		requests Action
	}{
		{
			name: "http to http request",
			services: Lifecycles(
				HTTPService(
					Name("myservice"),
					Port(12346),
					Proc(Name("echo"), EchoHandler()),
				),
			),
			requests: Actions(
				HTTPRequest(
					Port(12346),
					Body("test body"),
					Service("myservice"),
					Procedure("echo"),
					ExpectRespBody("test body"),
				),
			),
		},
		{
			name: "tchannel to tchannel request",
			services: Lifecycles(
				TChannelService(
					Name("myservice"),
					Port(22345),
					Proc(Name("echo"), EchoHandler()),
				),
			),
			requests: Actions(
				TChannelRequest(
					Port(22345),
					Body("test body"),
					Service("myservice"),
					Procedure("echo"),
					ExpectRespBody("test body"),
				),
			),
		},
		{
			name: "grpc to grpc request",
			services: Lifecycles(
				GRPCService(
					Name("myservice"),
					Port(32345),
					Proc(Name("echo"), EchoHandler()),
				),
			),
			requests: Actions(
				GRPCRequest(
					Port(32345),
					Body("test body"),
					Service("myservice"),
					Procedure("echo"),
					ExpectRespBody("test body"),
				),
			),
		},
		{
			name: "response errors",
			services: Lifecycles(
				HTTPService(
					Name("myservice"),
					Port(12346),
					Proc(Name("error"), ErrorHandler(errors.New("error from myservice"))),
				),
				TChannelService(
					Name("myotherservice"),
					Port(22346),
					Proc(Name("error"), ErrorHandler(errors.New("error from myotherservice"))),
				),
				GRPCService(
					Name("myotherservice2"),
					Port(32346),
					Proc(Name("error"), ErrorHandler(errors.New("error from myotherservice2"))),
				),
			),
			requests: Actions(
				HTTPRequest(
					Port(12346),
					Service("myservice"),
					Procedure("error"),
					ExpectError("error from myservice"),
				),
				TChannelRequest(
					Port(22346),
					Service("myotherservice"),
					Procedure("error"),
					ExpectError("error from myotherservice"),
				),
				GRPCRequest(
					Port(32346),
					Service("myotherservice2"),
					Procedure("error"),
					ExpectError("error from myotherservice2"),
				),
			),
		},
		{
			name: "ordered requests",
			services: Lifecycles(
				HTTPService(
					Name("myservice"),
					Port(11112),
					Proc(
						Name("proc"),
						OrderedRequestHandler(
							ErrorHandler(yarpcerrors.InternalErrorf("internal error")),
							StaticHandler("success"),
							EchoHandlerWithPrefix("echo: "),
							EchoHandler(),
						),
					),
				),
			),
			requests: Actions(
				HTTPRequest(
					Port(11112),
					Service("myservice"),
					Procedure("proc"),
					ShardKey("ignoreme"),
					ExpectError(yarpcerrors.InternalErrorf("internal error").Error()),
				),
				HTTPRequest(
					Port(11112),
					Service("myservice"),
					Procedure("proc"),
					ExpectRespBody("success"),
				),
				HTTPRequest(
					Port(11112),
					Service("myservice"),
					Procedure("proc"),
					Body("hello"),
					ExpectRespBody("echo: hello"),
				),
				HTTPRequest(
					Port(11112),
					Service("myservice"),
					Procedure("proc"),
					GiveAndExpectLargeBodyIsEchoed(1<<17),
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
