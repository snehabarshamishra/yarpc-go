package api

import "go.uber.org/yarpc/api/transport"

// ServerStreamAction is an action applied to a ServerStream.
// If the action returns an error, that error will be used to end the ServerStream.
type ServerStreamAction interface {
	Lifecycle

	ApplyServerStream(transport.ServerStream) error
}

// ServerStreamAction converts a function into a StreamAction.
type ServerStreamActionFunc func(transport.ServerStream) error

// ApplyServerStream implements ServerStreamAction.
func (f ServerStreamActionFunc) ApplyServerStream(c transport.ServerStream) error { return f(c) }

// Start is a noop for wrapped functions
func (f ServerStreamActionFunc) Start(TestingT) error { return nil }

// Stop is a noop for wrapped functions
func (f ServerStreamActionFunc) Stop(TestingT) error { return nil }

