package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
)

const (
	ContentType       = "Content-Type"
	ApplicationJSON   = "application/json"
	RabbitMQHaMode    = "ha-mode"
	RabbitMQShardMode = "sharding-mode"

	Success = "Success"
)

type (
	// QueueBody ...
	QueueBody struct {
		AutoDelete bool      `json:"auto_delete"`
		Durable    bool      `json:"durable"`
		Arguments  QueueAgmt `json:"arguments"`
	}

	// QueueAgmt ...
	QueueAgmt struct {
		MaxLength            int    `json:"x-max-length,omitempty"`
		MessageTTL           int    `json:"x-message-ttl,omitempty"`
		DeadLetterExchange   string `json:"x-dead-letter-exchange,omitempty"`
		DeadLetterRoutingKey string `json:"x-dead-letter-routing-key,omitempty"`
	}

	//BindingSources ...
	BindingSources struct {
		Destination     string `json:"destination"`
		DestinationType string `json:"destination_type"`
	}

	// Permission define request body of create policy api
	Permission struct {
		Pattern    string     `json:"pattern"`
		Definition Definition `json:"definition"`
		ApplyTo    string     `json:"apply-to"`
	}

	// Definition define parameters of policy
	Definition struct {
		HaMode        string `json:"ha-mode,omitempty"`
		HaParam       int    `json:"ha-params,omitempty"`
		HaSyncMode    string `json:"ha-sync-mode,omitempty"`
		ShardsPerNode int    `json:"shards-per-node,omitempty"`
	}

	// BindBody ...
	BindBody struct {
		RouteingKey string `json:"routing_key"`
	}

	// Queue ...
	Queue struct {
		Name                      string   `json:"name"`
		Consumers                 int      `json:"consumers"`
		ConsumersNum              int      `json:"number_of_consumers"`
		MessageNum                int      `json:"messages"`
		MessageSize               int      `json:"message_bytes"`
		Memory                    int      `json:"memory"`
		State                     string   `json:"state"`
		Node                      string   `json:"node"`
		SlaveNodes                []string `json:"slave_nodes"`
		Policy                    string   `json:"policy"`
		EffectivePolicyDefinition struct {
			HAParams int `json:"ha-params"`
		} `json:"effective_policy_definition"`
	}

	Metadata struct {
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

	MessageAttribute struct {
		Name  string `json:"name"`
		Value struct {
			DataType    string `json:"dataType"`
			StringValue string `json:"stringValue"`
		} `json:"value"`
	}
)

func parseMsgQueueName(name string) (accountID, queueName string) {
	msgq := strings.Split(name, "@")
	if len(msgq) < 2 {
		return
	}

	return msgq[0], msgq[1]
}

func parseQueueStruct(queue *Queue) (comQueue *common.Queue) {
	comQueue = &common.Queue{
		Name:         queue.Name,
		Consumers:    queue.Consumers,
		ConsumersNum: queue.ConsumersNum,
		MessageNum:   queue.MessageNum,
		MessageSize:  queue.MessageSize,
		Memory:       queue.Memory,
		State:        queue.State,
		Node:         queue.Node,
		SlaveNodes:   queue.SlaveNodes,
		Policy:       queue.Policy,
		ReplicaSize:  queue.EffectivePolicyDefinition.HAParams,
	}

	return
}

func queuesInfo2Map(queues []*common.Queue) (maps map[string]*common.Queue) {
	maps = map[string]*common.Queue{}
	for _, q := range queues {
		maps[q.Name] = q
	}
	return
}

// GetURL ...
func (q *QueueMethod) getURL(endpoint string) string {
	return fmt.Sprintf(
		"http://%s:%s@%s:%d/api/%s",
		q.Provider.AMQP.Account,
		q.Provider.AMQP.Password,
		q.Provider.AMQP.IP,
		q.Provider.AMQP.HTTPPort,
		endpoint,
	)
}
