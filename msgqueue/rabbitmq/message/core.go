package message

import (
	"github.com/streadway/amqp"
)

// Config define rabbitmq message config
type Config struct{}

// Checking ...
type Checking func(messageIDs string) (deleted bool, iMetadata []byte)

// Receiving ...
type Receiving func(d amqp.Delivery, iMetadata []byte)

type Metadata struct {
	VersionOfStruct        string             `json:"versionOfStruct"`
	MessageID              string             `json:"messageId"`
	UserID                 string             `json:"userId"`
	GroupID                string             `json:"groupId"`
	QueueName              string             `json:"queueName"`
	IsEncrypted            bool               `json:"isEncrypted"`
	KMSID                  string             `json:"kmsId"`
	SendTimestamp          string             `json:"sendTimestamp"`
	DisplayName            string             `json:"displayName"`
	MessageAttributes      []MessageAttribute `json:"messageAttributes"`
	MD5OfMessageAttributes string             `json:"md5OfMessageAttributes"`
	MD5OfMessageBody       string             `json:"md5OfMessageBody"`
}

type MessageAttribute struct {
	Name  string `json:"name"`
	Value struct {
		DataType    string `json:"dataType"`
		StringValue string `json:"stringValue"`
	} `json:"value"`
}
