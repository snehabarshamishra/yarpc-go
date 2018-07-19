package chooserbenchmark

import (
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/yarpc/api/transport"
)

func TestFewestPendingChooser(t *testing.T) {
	f, err := NewPeerListChooserFactory()
	assert.NoError(t, err)

	clientGroup := &ClientGroup{
		ListType:        FewestPending,
		ListUpdaterType: Static,
	}
	plc, err := f.CreatePeerListChooser(clientGroup, 2)
	assert.NoError(t, err)
	err = plc.Start()
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		ctx := context.Background()
		plc.Choose(ctx, &transport.Request{})

	}
}
