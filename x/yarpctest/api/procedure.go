package api

import "go.uber.org/yarpc/api/transport"

// ProcOpts are configuration options for a procedure.
type ProcOpts struct {
	Name        string
	HandlerSpec transport.HandlerSpec
}

// ProcOption defines the options that can be applied to a procedure.
type ProcOption interface {
	Lifecycle

	ApplyProc(*ProcOpts)
}
