package common

type (
	// RPCInterface ...
	RPCInterface interface {
		SendRPC(name, rKey, replyTo, payload string) (res string, err error)
		ReceiveRPC(name string, rec Receiving) (err error)
	}
)
