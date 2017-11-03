package types

import (
	"bytes"
	"io/ioutil"

	"github.com/stretchr/testify/require"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/x/yarpctest/api"
)

// SendStreamMsg is an action to send a message to a stream.  It can be
// applied to either a server or client stream.
type SendStreamMsg struct {
	api.TestingTInjectable
	api.NoopStop

	Msg         string
	WantErrMsgs []string
}

// ApplyClientStream implements ClientStreamAction
func (s *SendStreamMsg) ApplyClientStream(t api.TestingT, c transport.ClientStream) {
	_ = s.applyStream(t, c)
}

// ApplyServerStream implements ServerStreamAction
func (s *SendStreamMsg) ApplyServerStream(c transport.ServerStream) error {
	return s.applyStream(s.GetTestingT(), c)
}

func (s *SendStreamMsg) applyStream(t api.TestingT, c transport.BaseStream) error {
	err := c.SendMsg(&transport.StreamMessage{
		ReadCloser: ioutil.NopCloser(bytes.NewBufferString(s.Msg)),
	})
	if len(s.WantErrMsgs) > 0 {
		for _, wantErrMsg := range s.WantErrMsgs {
			require.Contains(t, err.Error(), wantErrMsg)
		}
		return err
	}
	require.NoError(t, err)
	return nil
}

// RecvStreamMsg is an action to receive a message from a stream.  It can
// be applied to either a server or client stream.
type RecvStreamMsg struct {
	api.TestingTInjectable
	api.NoopStop

	Msg         string
	WantErrMsgs []string
}

// ApplyClientStream implements ClientStreamAction
func (s *RecvStreamMsg) ApplyClientStream(t api.TestingT, c transport.ClientStream) {
	_ = s.applyStream(t, c)
}

// ApplyServerStream implements ServerStreamAction
func (s *RecvStreamMsg) ApplyServerStream(c transport.ServerStream) error {
	return s.applyStream(s.GetTestingT(), c)
}

func (s *RecvStreamMsg) applyStream(t api.TestingT, c transport.BaseStream) error {
	msg, err := c.RecvMsg()
	if len(s.WantErrMsgs) > 0 {
		require.Error(t, err)
		for _, wantErrMsg := range s.WantErrMsgs {
			require.Contains(t, err.Error(), wantErrMsg)
		}
		return err
	}
	require.NoError(t, err)

	actualMsg, err := ioutil.ReadAll(msg)
	require.NoError(t, err)
	require.Equal(t, bytes.NewBufferString(s.Msg).Bytes(), actualMsg, "mismatch on stream messages")
	return nil
}
