package common

type (
	MessageInterface interface {
		SendMessage(name, rKey, mid, payload string, headers map[string]interface{}) (err error)
		ReceiveMessage(name string, tgtsize int, del Checking, rec Receiving) (err error)
	}
	// Checking ...
	Checking func(messageIDs string) (deleted bool, iMetadata []byte)

	// Receiving ...
	Receiving func(messageIDs string, headers map[string]interface{}, body, iMetadata []byte)
)
