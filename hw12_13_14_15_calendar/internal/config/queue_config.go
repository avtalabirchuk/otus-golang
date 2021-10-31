package config

type QueueConfig struct {
	Host                 string `yaml:"host" config:"required"`
	Port                 int    `yaml:"port" config:"required"`
	User                 string `yaml:"user" config:"required"`
	Pass                 string `yaml:"pass" config:"required"`
	QueueName            string `yaml:"queueName" config:"required"`
	ExchangeType         string `yaml:"exchangeType" config:"required"`
	ScanTimeoutMs        int    `yaml:"scanTimeoutMs" config:"required"`
	QosPrefetchCount     int    `yaml:"qosPrefetchCount" config:"required"`
	MaxReconnectAttempts int    `yaml:"maxReconnectAttempts" config:"required"`
	ReconnectTimeoutMs   int
}

func defaultQueueConfig() QueueConfig {
	return QueueConfig{
		Host:                 "localhost",
		Port:                 5672,
		QueueName:            "events",
		ExchangeType:         "direct",
		ScanTimeoutMs:        10000,
		QosPrefetchCount:     50,
		MaxReconnectAttempts: 5,
		ReconnectTimeoutMs:   2000,
	}
}
