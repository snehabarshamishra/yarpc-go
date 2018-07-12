package loadbalancingbenchmark

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/yarpc/api/peer"
	"go.uber.org/yarpc/peer/pendingheap"
	"go.uber.org/yarpc/peer/roundrobin"
)

type ListFactory func(group *ClientGroup) (peer.ChooserList, error)

var listFactories = make(map[ListType]ListFactory)

func NewRoundRobinList(group *ClientGroup) (peer.ChooserList, error) {
	return roundrobin.New(NewBenchTransport()), nil
}

func NewFewestPendingList(group *ClientGroup) (peer.ChooserList, error) {
	return pendingheap.New(NewBenchTransport()), nil
}

func RegisterList(balancingType ListType, factory ListFactory) error {
	if factory == nil {
		return errors.New(fmt.Sprintf(`unable to register %q, factory is nil`, balancingType))
	}
	if _, ok := listFactories[balancingType]; ok {
		return errors.New(fmt.Sprintf(`factory for %q already exists`, balancingType))
	}
	listFactories[balancingType] = factory
	return nil
}

func InitListFactory() error {
	return multierr.Combine(
		RegisterList(RoundRobin, NewRoundRobinList),
		RegisterList(FewestPending, NewFewestPendingList),
	)
}

func CreateList(group *ClientGroup) (peer.ChooserList, error) {
	factory, ok := listFactories[group.LType]
	if !ok {
		return nil, errors.New(fmt.Sprintf(`type %q is not supported`, group.LType))
	}
	return factory(group)
}
