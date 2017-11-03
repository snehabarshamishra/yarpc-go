package api

import (
	"bytes"

	"go.uber.org/yarpc/api/transport"
)


// RequestOpts are configuration options for a yarpc Request and assertions
// to make on the response.
type RequestOpts struct {
	Port             int
	GiveRequest      *transport.Request
	ExpectedResponse *transport.Response
	ExpectedError    error
}

// NewRequestOpts initializes a RequestOpts struct.
func NewRequestOpts() RequestOpts {
	return RequestOpts{
		GiveRequest: &transport.Request{
			Caller:   "unknown",
			Encoding: transport.Encoding("raw"),
			Body:     bytes.NewBufferString(""),
		},
		ExpectedResponse: &transport.Response{},
	}
}

// RequestOption can be used to configure a request.
type RequestOption interface {
	ApplyRequest(*RequestOpts)
}

// RequestOptionFunc converts a function into a RequestOption.
type RequestOptionFunc func(*RequestOpts)

// ApplyRequest implements RequestOption.
func (f RequestOptionFunc) ApplyRequest(opts *RequestOpts) { f(opts) }
