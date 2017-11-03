package yarpctest

import (
	"go.uber.org/yarpc/x/yarpctest/types"
)

// SendStreamMsg sends a message to a stream.
func SendStreamMsg(sendMsg string) *types.SendStreamMsg {
	return &types.SendStreamMsg{Msg: sendMsg}
}

// SendStreamMsgAndExpectError sends a message on a stream and asserts on the
// error returned.
func SendStreamMsgAndExpectError(sendMsg string, wantErrMsgs ...string) *types.SendStreamMsg {
	return &types.SendStreamMsg{Msg: sendMsg, WantErrMsgs: wantErrMsgs}
}

// RecvStreamMsg waits to receive a message on a client stream.
func RecvStreamMsg(wantMsg string) *types.RecvStreamMsg {
	return &types.RecvStreamMsg{Msg: wantMsg}
}

// RecvStreamErr waits to receive a message on a client stream.  It expects
// an error.
func RecvStreamErr(wantErrMsgs ...string) *types.RecvStreamMsg {
	return &types.RecvStreamMsg{WantErrMsgs: wantErrMsgs}
}
