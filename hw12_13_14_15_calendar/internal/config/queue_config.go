package config

type QueueConfig struct {
	URI                  string `yaml:"uri" env:"QUEUE_URI"`
	QueueName            string `yaml:"queueName" env:"QUEUE_NAME" env-default:"events"`
	ExchangeType         string `yaml:"exchangeType" env:"QUEUE_EXCHANGE_TYPE" env-default:"direct"`
	ScanTimeoutMs        int    `yaml:"scanTimeoutMs" env:"QUEUE_SCAN_TIMEOUT_MS" env-default:"10000"`
	QosPrefetchCount     int    `yaml:"qosPrefetchCount" env:"QUEUE_QOS_PREFETCH_COUNT" env-default:"50"`
	MaxReconnectAttempts int    `yaml:"maxReconnectAttempts" env:"QUEUE_MAC_RECONNECT_ATTEMPTS" env-default:"5"`
	ReconnectTimeoutMs   int    `yaml:"reconnectTimeoutMs" env:"QUEUE_MAC_RECONNECT_TIMEOUT_MS" env-default:"2000"`
}
