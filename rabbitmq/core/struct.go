package core

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
)
