package api

import (
	"context"

	"go.uber.org/yarpc/api/transport"
)

// UnaryHandlerFunc converts a function into a transport.UnaryHandler.
type UnaryHandlerFunc func(context.Context, *transport.Request, transport.ResponseWriter) error

// Handle implements yarpc/api/transport#UnaryHandler.
func (f UnaryHandlerFunc) Handle(ctx context.Context, req *transport.Request, resw transport.ResponseWriter) error {
	return f(ctx, req, resw)
}

// HandlerOpts are configuration options for a series of handlers.
type HandlerOpts struct {
	Handlers []transport.UnaryHandler
}

// HandlerOption defines options that can be passed into a handler.
type HandlerOption interface {
	Lifecycle

	ApplyHandler(opts *HandlerOpts)
}

// HandlerOptionFunc converts a function into a HandlerOption.
type HandlerOptionFunc func(opts *HandlerOpts)

// ApplyHandler implements HandlerOption.
func (f HandlerOptionFunc) ApplyHandler(opts *HandlerOpts) { f(opts) }
