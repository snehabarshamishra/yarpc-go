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
	"context"
	"fmt"
	"go.uber.org/atomic"
	"go.uber.org/yarpc/api/transport"
	"io"
	"io/ioutil"
)

// ProcOpts are configuration options for a procedure.
type ProcOpts struct {
	Name        string
	HandlerSpec transport.HandlerSpec
}

// ProcOption defines the options that can be applied to a procedure.
type ProcOption interface {
	ApplyProc(*ProcOpts)
}

// ProcOptionFunc converts a function into a ProcOption.
type ProcOptionFunc func(opts *ProcOpts)

// ApplyProc implements ProcOption.
func (f ProcOptionFunc) ApplyProc(opts *ProcOpts) { f(opts) }

// Proc will create a new Procedure that can be included in a Service.
func Proc(options ...ProcOption) ServiceOption {
	return ServiceOptionFunc(func(opts *ServiceOpts) {
		opts.Procedures = append(opts.Procedures, newProc(options...))
	})
}

func newProc(options ...ProcOption) transport.Procedure {
	opts := ProcOpts{Name: "proc"}
	for _, option := range options {
		option.ApplyProc(&opts)
	}
	return transport.Procedure{
		Name:        opts.Name,
		HandlerSpec: opts.HandlerSpec,
	}
}

// UnaryHandlerFunc converts a function into a UnaryHandler.
type UnaryHandlerFunc func(context.Context, *transport.Request, transport.ResponseWriter) error

// Handle implements yarpc/api/transport#UnaryHandler.
func (f UnaryHandlerFunc) Handle(ctx context.Context, req *transport.Request, resw transport.ResponseWriter) error {
	return f(ctx, req, resw)
}

// UnaryHandler is a struct that implements the ProcOptions and HandlerOption
// interfaces (so it can be used directly as a procedure, or as a single use
// handler (depending on the use case)).
type UnaryHandler struct {
	H transport.UnaryHandler
}

// ApplyProc implements ProcOption.
func (h *UnaryHandler) ApplyProc(opts *ProcOpts) {
	opts.HandlerSpec = transport.NewUnaryHandlerSpec(h.H)
}

// ApplyHandler implements HandlerOption.
func (h *UnaryHandler) ApplyHandler(opts *HandlerOpts) {
	opts.Handlers = append(opts.Handlers, h.H)
}

// EchoHandler is a Unary Handler that will echo the body of the request
// into the response.
func EchoHandler() *UnaryHandler {
	return &UnaryHandler{H: newEchoHandler()}
}

func newEchoHandler() transport.UnaryHandler {
	return UnaryHandlerFunc(
		func(_ context.Context, req *transport.Request, resw transport.ResponseWriter) error {
			_, err := io.Copy(resw, req.Body)
			return err
		},
	)
}

// StaticHandler will always return the same response.
func StaticHandler(msg string) *UnaryHandler {
	return &UnaryHandler{H: newStaticHandler(msg)}
}

func newStaticHandler(msg string) transport.UnaryHandler {
	return UnaryHandlerFunc(
		func(_ context.Context, _ *transport.Request, resw transport.ResponseWriter) error {
			_, err := io.WriteString(resw, msg)
			return err
		},
	)
}

// ErrorHandler will always return an Error.
func ErrorHandler(err error) *UnaryHandler {
	return &UnaryHandler{H: newErrorHandler(err)}
}

func newErrorHandler(err error) transport.UnaryHandler {
	return UnaryHandlerFunc(
		func(context.Context, *transport.Request, transport.ResponseWriter) error {
			return err
		},
	)
}

// EchoHandlerWithPrefix will echo the request it receives into the
// response, but, it will insert a prefix in front of the message.
func EchoHandlerWithPrefix(prefix string) *UnaryHandler {
	return &UnaryHandler{H: newEchoHandlerWithPrefix(prefix)}
}

func newEchoHandlerWithPrefix(prefix string) transport.UnaryHandler {
	return UnaryHandlerFunc(
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

// HandlerOpts are configuration options for a series of handlers.
type HandlerOpts struct {
	Handlers []transport.UnaryHandler
}

// HandlerOption defines options that can be passed into a handler.
type HandlerOption interface {
	ApplyHandler(opts *HandlerOpts)
}

// HandlerOptionFunc converts a function into a HandlerOption.
type HandlerOptionFunc func(opts *HandlerOpts)

// ApplyHandler implements HandlerOption.
func (f HandlerOptionFunc) ApplyHandler(opts *HandlerOpts) { f(opts) }

// OrderedRequestHandler will execute a series of Handlers in the order they
// are passed in.  If the number of requests does not match, it will return an
// unknown error.
func OrderedRequestHandler(options ...HandlerOption) *OrderedHandler {
	opts := HandlerOpts{}
	for _, option := range options {
		option.ApplyHandler(&opts)
	}
	return &OrderedHandler{
		handlers: opts.Handlers,
	}
}

// OrderedHandler implements the transport.UnaryHandler and ProcOption
// interfaces.
type OrderedHandler struct {
	attempt  atomic.Int32
	handlers []transport.UnaryHandler
}

// ApplyProc implements ProcOption.
func (h *OrderedHandler) ApplyProc(opts *ProcOpts) {
	opts.HandlerSpec = transport.NewUnaryHandlerSpec(h)
}

// Handle implements transport.UnaryHandler#Handle.
func (h *OrderedHandler) Handle(ctx context.Context, req *transport.Request, resw transport.ResponseWriter) error {
	if len(h.handlers) <= 0 {
		return fmt.Errorf("No handlers for the request")
	}
	n := h.attempt.Inc()
	if int(n) > len(h.handlers) {
		return fmt.Errorf("too many requests, expected %d, got %d", len(h.handlers), n)
	}
	return h.handlers[n-1].Handle(ctx, req, resw)
}
