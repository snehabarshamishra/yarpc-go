package types

import (
	"context"
	"fmt"

	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/x/yarpctest/api"
)

// UnaryHandler is a struct that implements the ProcOptions and HandlerOption
// interfaces (so it can be used directly as a procedure, or as a single use
// handler (depending on the use case)).
type UnaryHandler struct {
	api.TestingTInjectable
	api.NoopStop

	H transport.UnaryHandler
}

// ApplyProc implements ProcOption.
func (h *UnaryHandler) ApplyProc(opts *api.ProcOpts) {
	opts.HandlerSpec = transport.NewUnaryHandlerSpec(h.H)
}

// ApplyHandler implements HandlerOption.
func (h *UnaryHandler) ApplyHandler(opts *api.HandlerOpts) {
	opts.Handlers = append(opts.Handlers, h.H)
}

// OrderedHandler implements the transport.UnaryHandler and ProcOption
// interfaces.
type OrderedHandler struct {
	api.TestingTInjectable
	api.NoopStop

	attempt  atomic.Int32
	Handlers []transport.UnaryHandler
}

// ApplyProc implements ProcOption.
func (h *OrderedHandler) ApplyProc(opts *api.ProcOpts) {
	opts.HandlerSpec = transport.NewUnaryHandlerSpec(h)
}

// Handle implements transport.UnaryHandler#Handle.
func (h *OrderedHandler) Handle(ctx context.Context, req *transport.Request, resw transport.ResponseWriter) error {
	if len(h.Handlers) <= 0 {
		return fmt.Errorf("No handlers for the request")
	}
	n := h.attempt.Inc()
	if int(n) > len(h.Handlers) {
		return fmt.Errorf("too many requests, expected %d, got %d", len(h.Handlers), n)
	}
	return h.Handlers[n-1].Handle(ctx, req, resw)
}
