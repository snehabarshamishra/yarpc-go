package api

import "go.uber.org/yarpc/api/transport"

// ClientStreamRequestOpts are configuration options for a yarpc stream request.
type ClientStreamRequestOpts struct {
	Port            int
	GiveRequestMeta *transport.RequestMeta
	StreamActions   []ClientStreamAction
}

// NewClientStreamRequestOpts initializes a ClientStreamRequestOpts struct.
func NewClientStreamRequestOpts() ClientStreamRequestOpts {
	return ClientStreamRequestOpts{
		GiveRequestMeta: &transport.RequestMeta{
			Caller:   "unknown",
			Encoding: transport.Encoding("raw"),
		},
	}
}

// ClientStreamRequestOption can be used to configure a request.
type ClientStreamRequestOption interface {
	ApplyClientStreamRequest(*ClientStreamRequestOpts)
}

// ClientStreamRequestOptionFunc converts a function into a ClientStreamRequestOption.
type ClientStreamRequestOptionFunc func(*ClientStreamRequestOpts)

// ApplyClientStreamRequest implements ClientStreamRequestOption.
func (f ClientStreamRequestOptionFunc) ApplyClientStreamRequest(opts *ClientStreamRequestOpts) { f(opts) }

// ClientStreamAction is an action applied to a ClientStream.
type ClientStreamAction interface {
	ApplyClientStream(TestingT, transport.ClientStream)
}

// ClientStreamAction converts a function into a StreamAction.
type ClientStreamActionFunc func(TestingT, transport.ClientStream)

// ApplyClientStream implements ClientStreamAction.
func (f ClientStreamActionFunc) ApplyClientStream(t TestingT, c transport.ClientStream) { f(t, c) }
