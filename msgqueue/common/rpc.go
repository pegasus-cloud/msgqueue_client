package common

type (
	// RPCInterface ...
	RPCInterface interface {
		SendRPC(name, rKey, replyTo, payload string) (res string, err error)
		ReceiveRPC(name string, del Delivering) (err error)
	}

	// Delivering return data from server side of rpc mode
	Delivering func(messageIDs string, headers map[string]interface{}, body []byte) (data []byte)
)
