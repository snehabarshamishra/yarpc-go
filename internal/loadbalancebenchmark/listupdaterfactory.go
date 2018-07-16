package loadbalancingbenchmark

import (
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/yarpc/api/peer"
	peer2 "go.uber.org/yarpc/peer"
)

type ListUpdaterFactory func(serverCount int) (peer.Binder, error)

var listUpdaterFactories = make(map[ListUpdaterType]ListUpdaterFactory)

func NewStaticListUpdater(serverCount int) (peer.Binder, error) {
	return peer2.BindPeers(MakePeerIdentifier(serverCount)), nil
}

func RegisterListUpdater(updaterType ListUpdaterType, factory ListUpdaterFactory) error {
	if factory == nil {
		return fmt.Errorf(`unable to register %q, factory is nil`, updaterType)
	}
	if _, ok := listUpdaterFactories[updaterType]; ok {
		return fmt.Errorf(`factory for %q already exists`, updaterType)
	}
	listUpdaterFactories[updaterType] = factory
	return nil
}

func InitListUpdaterFactory() error {
	return multierr.Combine(
		RegisterListUpdater(Static, NewStaticListUpdater),
	)
}

func CreateListUpdater(group *ClientGroup, serverCount int) (peer.Binder, error) {
	factory, ok := listUpdaterFactories[group.LUType]
	if !ok {
		return nil, fmt.Errorf(`type %q is not supported`, group.LUType)
	}
	return factory(serverCount)
}
