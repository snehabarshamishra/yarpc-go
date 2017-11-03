package yarpctest

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"go.uber.org/yarpc/x/yarpctest/api"
)

type Action api.Action
type Lifecycle api.Lifecycle

// Lifecycles is a wrapper around a list of Lifecycle definitions.
func Lifecycles(l ...api.Lifecycle) api.Lifecycle {
	return lifecycles(l)
}

type lifecycles []api.Lifecycle

// Start the lifecycles. If there are any errors, stop any started lifecycles
// and fail the test.
func (ls lifecycles) Start(t api.TestingT) error {
	startedLifecycles := make(lifecycles, 0, len(ls))
	for _, l := range ls {
		err := l.Start(t)
		if !assert.NoError(t, err) {
			// Cleanup started lifecycles (this could fail)
			return multierr.Append(err, startedLifecycles.Stop(t))
		}
		startedLifecycles = append(startedLifecycles, l)
	}
	return nil
}

// Stop the lifecycles. Record all errors. If any lifecycle failed to stop
// fail the test.
func (ls lifecycles) Stop(t api.TestingT) error {
	var err error
	for _, l := range ls {
		err = multierr.Append(err, l.Stop(t))
	}
	assert.NoError(t, err)
	return err
}

// Actions will wrap a list of actions in a sequential executor.
func Actions(actions ...api.Action) api.Action {
	return multi(actions)
}

type multi []api.Action

func (m multi) Run(t api.TestingT) {
	for i, req := range m {
		api.Run(fmt.Sprintf("Action #%d", i), t, req.Run)
	}
}

