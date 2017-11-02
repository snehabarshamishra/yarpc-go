package yarpctest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/transport/http"
	"go.uber.org/yarpc/transport/tchannel"
)

// Actions will wrap a list of actions in a sequential executor.
func Actions(actions ...Action) Action {
	return multi(actions)
}

type multi []Action

func (m multi) Run(t TestingT) {
	for i, req := range m {
		Run(fmt.Sprintf("Req #%d", i), t, req.Run)
	}
}

// Action defines an object that can be "Run" to assert things against the
// world through action.
type Action interface {
	Run(t TestingT)
}

// RequestFunc is a helper to convert a function to implement the Request
// interface
type RequestFunc func(TestingT)

// Run implement Action.
func (f RequestFunc) Run(t TestingT) { f(t) }

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

// HTTPRequest creates a new YARPC http request.
func HTTPRequest(options ...RequestOption) Action {
	opts := NewRequestOpts()
	for _, option := range options {
		option.ApplyRequest(&opts)
	}
	return RequestFunc(func(t TestingT) {
		trans := http.NewTransport()
		out := trans.NewSingleOutbound(fmt.Sprintf("http://127.0.0.1:%d/", opts.Port))

		require.NoError(t, trans.Start())
		defer func() { assert.NoError(t, trans.Stop()) }()

		require.NoError(t, out.Start())
		defer func() { assert.NoError(t, out.Stop()) }()

		resp, cancel, err := sendRequest(out, opts.GiveRequest)
		defer cancel()
		validateError(t, err, opts.ExpectedError)
		if opts.ExpectedError == nil {
			validateResponse(t, resp, opts.ExpectedResponse)
		}
	})
}

// TChannelRequest creates a new tchannel request.
func TChannelRequest(options ...RequestOption) Action {
	opts := NewRequestOpts()
	for _, option := range options {
		option.ApplyRequest(&opts)
	}
	return RequestFunc(func(t TestingT) {
		trans, err := tchannel.NewTransport(tchannel.ServiceName(opts.GiveRequest.Caller))
		require.NoError(t, err)
		out := trans.NewSingleOutbound(fmt.Sprintf("127.0.0.1:%d", opts.Port))

		require.NoError(t, trans.Start())
		defer func() { assert.NoError(t, trans.Stop()) }()

		require.NoError(t, out.Start())
		defer func() { assert.NoError(t, out.Stop()) }()

		resp, cancel, err := sendRequest(out, opts.GiveRequest)
		defer cancel()
		validateError(t, err, opts.ExpectedError)
		if opts.ExpectedError == nil {
			validateResponse(t, resp, opts.ExpectedResponse)
		}
	})
}

func sendRequest(out transport.UnaryOutbound, request *transport.Request) (*transport.Response, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := out.Call(ctx, request)
	return resp, cancel, err
}

func validateError(t TestingT, actualErr error, wantError error) {
	if wantError != nil {
		require.Error(t, actualErr)
		require.Contains(t, actualErr.Error(), wantError.Error())
		return
	}
	require.NoError(t, actualErr)
}

func validateResponse(t TestingT, actualResp *transport.Response, expectedResp *transport.Response) {
	var actualBody []byte
	var expectedBody []byte
	var err error
	if actualResp.Body != nil {
		actualBody, err = ioutil.ReadAll(actualResp.Body)
		require.NoError(t, err)
	}
	if expectedResp.Body != nil {
		expectedBody, err = ioutil.ReadAll(expectedResp.Body)
		require.NoError(t, err)
	}
	assert.Equal(t, string(actualBody), string(expectedBody))
}

// Body sets the body on a request to the raw representation of the msg field.
func Body(msg string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.GiveRequest.Body = bytes.NewBufferString(msg)
	})
}

// Service sets the request service field.
func Service(service string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.GiveRequest.Service = service
	})
}

// Procedure sets the request procedure field.
func Procedure(proc string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.GiveRequest.Procedure = proc
	})
}

// ShardKey sets the request shard key field.
func ShardKey(key string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.GiveRequest.ShardKey = key
	})
}

// ExpectError creates an assertion on the request response to validate the
// error.
func ExpectError(errMsg string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.ExpectedError = errors.New(errMsg)
	})
}

// ExpectRespBody will assert that the response body matches at the end of the
// request.
func ExpectRespBody(body string) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		opts.ExpectedResponse.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	})
}

// GiveAndExpectLargeBodyIsEchoed creates an extremely large random byte buffer
// and validates that the body is echoed back to the response.
func GiveAndExpectLargeBodyIsEchoed(numOfBytes int) RequestOption {
	return RequestOptionFunc(func(opts *RequestOpts) {
		body := bytes.Repeat([]byte("t"), numOfBytes)
		opts.GiveRequest.Body = bytes.NewReader(body)
		opts.ExpectedResponse.Body = ioutil.NopCloser(bytes.NewReader(body))
	})
}
