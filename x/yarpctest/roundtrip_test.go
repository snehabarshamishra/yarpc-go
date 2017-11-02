package yarpctest

import (
	"testing"

	"errors"
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
