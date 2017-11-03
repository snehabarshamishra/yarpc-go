package types

import "go.uber.org/yarpc/x/yarpctest/api"

// Port is a concrete type that implements multiple interfaces as ease
// of use for synchronizing ports.
type Port struct {
	api.NoopLifecycle

	Port int
}

// ApplyService implements api.ServiceOption.
func (n *Port) ApplyService(opts *api.ServiceOpts) {
	opts.Port = n.Port
}

// ApplyRequest implements RequestOption
func (n *Port) ApplyRequest(opts *api.RequestOpts) {
	opts.Port = n.Port
}

// ApplyClientStreamRequest implements ClientStreamRequestOption
func (n *Port) ApplyClientStreamRequest(opts *api.ClientStreamRequestOpts) {
	opts.Port = n.Port
}
