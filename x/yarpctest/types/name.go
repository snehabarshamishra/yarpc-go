package types

import "go.uber.org/yarpc/x/yarpctest/api"

// Name is a concrete type that implements both ServiceOption and
// ProcedureOption interfaces so it can be used interchangeably.
type Name struct {
	api.NoopLifecycle

	Name string
}

// ApplyService implements ServiceOption.
func (n *Name) ApplyService(opts *api.ServiceOpts) {
	opts.Name = n.Name
}

// ApplyProc implements ProcOption.
func (n *Name) ApplyProc(opts *api.ProcOpts) {
	opts.Name = n.Name
}
