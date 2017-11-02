package yarpctest

// Port is a shared option across services and requests. It can be embedded
// in the constructor of both types.
func Port(port int) *PortOption {
	return &PortOption{port: port}
}

// PortOption is a concrete type that implements both ServiceOption and
// RequestOption interfaces.
type PortOption struct {
	port int
}

// ApplyService implements ServiceOption.
func (n *PortOption) ApplyService(opts *ServiceOpts) {
	opts.Port = n.port
}

// ApplyRequest implements RequestOption
func (n *PortOption) ApplyRequest(opts *RequestOpts) {
	opts.Port = n.port
}
