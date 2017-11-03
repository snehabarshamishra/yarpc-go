package types

import "go.uber.org/yarpc/x/yarpctest/api"

// Service is a concrete type that represents the "service" for a request.
// It can be used in multiple interfaces.
type Service struct {
	Service string
}

// ApplyRequest implements api.RequestOption
func (n *Service) ApplyRequest(opts *api.RequestOpts) {
	opts.GiveRequest.Service = n.Service
}

// ApplyClientStreamRequest implements api.ClientStreamRequestOption
func (n *Service) ApplyClientStreamRequest(opts *api.ClientStreamRequestOpts) {
	opts.GiveRequestMeta.Service = n.Service
}

// Procedure is a concrete type that represents the "procedure" for a request.
// It can be used in multiple interfaces.
type Procedure struct {
	Procedure string
}

// ApplyRequest implements api.RequestOption
func (n *Procedure) ApplyRequest(opts *api.RequestOpts) {
	opts.GiveRequest.Procedure = n.Procedure
}

// ApplyClientStreamRequest implements api.ClientStreamRequestOption
func (n *Procedure) ApplyClientStreamRequest(opts *api.ClientStreamRequestOpts) {
	opts.GiveRequestMeta.Procedure = n.Procedure
}

// ShardKey is a concrete type that represents the "shard key" for a request.
// It can be used in multiple interfaces.
type ShardKey struct {
	ShardKey string
}

// ApplyRequest implements api.RequestOption
func (n *ShardKey) ApplyRequest(opts *api.RequestOpts) {
	opts.GiveRequest.ShardKey = n.ShardKey
}

// ApplyClientStreamRequest implements api.ClientStreamRequestOption
func (n *ShardKey) ApplyClientStreamRequest(opts *api.ClientStreamRequestOpts) {
	opts.GiveRequestMeta.ShardKey = n.ShardKey
}
