package yarpctest

import (
	"context"
	"io"
	"io/ioutil"

	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/x/yarpctest/api"
	"go.uber.org/yarpc/x/yarpctest/types"
)

// EchoHandler is a Unary Handler that will echo the body of the request
// into the response.
func EchoHandler() *types.UnaryHandler {
	return &types.UnaryHandler{H: newEchoHandler()}
}

func newEchoHandler() transport.UnaryHandler {
	return api.UnaryHandlerFunc(
		func(_ context.Context, req *transport.Request, resw transport.ResponseWriter) error {
			_, err := io.Copy(resw, req.Body)
			return err
		},
	)
}

// StaticHandler will always return the same response.
func StaticHandler(msg string) *types.UnaryHandler {
	return &types.UnaryHandler{H: newStaticHandler(msg)}
}

func newStaticHandler(msg string) transport.UnaryHandler {
	return api.UnaryHandlerFunc(
		func(_ context.Context, _ *transport.Request, resw transport.ResponseWriter) error {
			_, err := io.WriteString(resw, msg)
			return err
		},
	)
}

// ErrorHandler will always return an Error.
func ErrorHandler(err error) *types.UnaryHandler {
	return &types.UnaryHandler{H: newErrorHandler(err)}
}

func newErrorHandler(err error) transport.UnaryHandler {
	return api.UnaryHandlerFunc(
		func(context.Context, *transport.Request, transport.ResponseWriter) error {
			return err
		},
	)
}

// EchoHandlerWithPrefix will echo the request it receives into the
// response, but, it will insert a prefix in front of the message.
func EchoHandlerWithPrefix(prefix string) *types.UnaryHandler {
	return &types.UnaryHandler{H: newEchoHandlerWithPrefix(prefix)}
}

func newEchoHandlerWithPrefix(prefix string) transport.UnaryHandler {
	return api.UnaryHandlerFunc(
		func(_ context.Context, req *transport.Request, resw transport.ResponseWriter) error {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}
			newMsg := prefix + string(body)
			_, err = io.WriteString(resw, newMsg)
			return err
		},
	)
}

// OrderedRequestHandler will execute a series of Handlers in the order they
// are passed in.  If the number of requests does not match, it will return an
// unknown error.
func OrderedRequestHandler(options ...api.HandlerOption) *types.OrderedHandler {
	opts := api.HandlerOpts{}
	for _, option := range options {
		option.ApplyHandler(&opts)
	}
	return &types.OrderedHandler{
		Handlers: opts.Handlers,
	}
}
