package api

import "go.uber.org/yarpc/api/transport"

// ServiceOpts are the configuration options for a yarpc service.
type ServiceOpts struct {
	Name       string
	Port       int
	Procedures []transport.Procedure
}

// ServiceOption is an option when creating a Service.
type ServiceOption interface {
	Lifecycle

	ApplyService(*ServiceOpts)
}

// ServiceOptionFunc converts a function into a ServiceOption.
type ServiceOptionFunc func(*ServiceOpts)

// ApplyService implements ServiceOption.
func (f ServiceOptionFunc) ApplyService(opts *ServiceOpts) { f(opts) }

// Start is a noop for wrapped functions
func (f ServiceOptionFunc) Start(TestingT) error { return nil }

// Stop is a noop for wrapped functions
func (f ServiceOptionFunc) Stop(TestingT) error { return nil }
