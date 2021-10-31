package config

type QueueConfig struct {
	Host             string `yaml:"host" config:"required"`
	Port             int    `yaml:"port" config:"required"`
	User             string `yaml:"user" config:"required"`
	Pass             string `yaml:"pass" config:"required"`
	QueueName        string `yaml:"queueName" config:"required"`
	ExchangeType     string `yaml:"exchangeType" config:"required"`
	ScanTimeout      int    `yaml:"scanTimeout" config:"required"`
	QosPrefetchCount int    `yaml:"qosPrefetchCount" config:"required"`
}

func defaultQueueConfig() QueueConfig {
	return QueueConfig{
		Host:             "localhost",
		Port:             5672,
		QueueName:        "events",
		ExchangeType:     "direct",
		ScanTimeout:      10,
		QosPrefetchCount: 50,
	}
}
