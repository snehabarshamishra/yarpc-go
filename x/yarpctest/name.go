package yarpctest

// Name is a shared option across services and procedures. It can be used as
// input in the constructor of both types.
func Name(name string) *NameOption {
	return &NameOption{name: name}
}

// NameOption is a concrete type that implements both ServiceOption and
// ProcedureOption interfaces.
type NameOption struct {
	name string
}

// ApplyService implements ServiceOption.
func (n *NameOption) ApplyService(opts *ServiceOpts) {
	opts.Name = n.name
}

// ApplyProc implements ProcOption.
func (n *NameOption) ApplyProc(opts *ProcOpts) {
	opts.Name = n.name
}
