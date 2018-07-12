package loadbalancingbenchmark

import (
	"go.uber.org/multierr"
	"go.uber.org/yarpc/peer"
)

func InitFactory() error {
	return multierr.Combine(
		InitListFactory(),
		InitListUpdaterFactory(),
	)
}

func CreatePeerListChooser(group *ClientGroup, serverCount int) (*peer.BoundChooser, error) {
	list, err := CreateList(group)
	if err != nil {
		return nil, err
	}
	listUpdater, err := CreateListUpdater(group, serverCount)
	if err != nil {
		return nil, err
	}
	return peer.Bind(list, listUpdater), nil
}
