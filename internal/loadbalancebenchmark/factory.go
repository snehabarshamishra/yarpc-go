package loadbalancingbenchmark

import (
	"errors"
	"fmt"
	"go.uber.org/yarpc/peer"
	"go.uber.org/yarpc/peer/pendingheap"
	"go.uber.org/yarpc/peer/roundrobin"
)

type PeerChooserFactory func(config *TestConfig) (*peer.BoundChooser, error)

var peerChooserFactories = make(map[ChooserType]PeerChooserFactory)

func NewRoundRobinPeerChooser(config *TestConfig) (*peer.BoundChooser, error) {
	transport := NewBenchTransport()
	list := roundrobin.New(transport)
	// TODO, add support for non-static peer list updater, but now just start from static
	cnt, err := config.GetServerCount()
	if err != nil {
		return nil, err
	}
	return peer.Bind(list, peer.BindPeers(MakePeerIdentifier(cnt))), nil
}

func NewFewestPendingPeerChooser(config *TestConfig) (*peer.BoundChooser, error) {
	transport := NewBenchTransport()
	list := pendingheap.New(transport)
	// TODO, add support for non-static peer list updater, but now just start from static
	cnt, err := config.GetServerCount()
	if err != nil {
		return nil, err
	}
	return peer.Bind(list, peer.BindPeers(MakePeerIdentifier(cnt))), nil
}

func Register(chooserType ChooserType, factory PeerChooserFactory) error {
	if factory == nil {
		return errors.New(fmt.Sprintf(`unable to register %q, factory is nil`, chooserType))
	}
	if _, ok := peerChooserFactories[chooserType]; ok {
		return errors.New(fmt.Sprintf(`factory for %q already exists`, chooserType))
	}
	peerChooserFactories[chooserType] = factory
	return nil
}

func Init() {
	Register(RoundRobin, NewRoundRobinPeerChooser)
	Register(FewestPending, NewFewestPendingPeerChooser)
}

func CreatePeerChooser(config *TestConfig) (*peer.BoundChooser, error) {
	chooserType, err := config.GetChooserType()
	if err != nil {
		return nil, err
	}
	peerChooserFactory, ok := peerChooserFactories[chooserType]
	if !ok {
		return nil, errors.New(fmt.Sprintf(`ChooserType %q is not supported`, chooserType))
	}
	return peerChooserFactory(config)
}
