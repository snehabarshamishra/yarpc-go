package api

// Lifecycle defines test infra that needs to be started before the actions
// and stopped afterwards.
type Lifecycle interface {
	Start(TestingT) error
	Stop(TestingT) error
}

// Action defines an object that can be "Run" to assert things against the
// world through action.
type Action interface {
	Run(t TestingT)
}

// ActionFunc is a helper to convert a function to implement the Action
// interface
type ActionFunc func(TestingT)

// Run implement Action.
func (f ActionFunc) Run(t TestingT) { f(t) }
