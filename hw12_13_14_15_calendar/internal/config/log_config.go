package config

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

func defaultLogConfig() LogConfig {
	return LogConfig{
		Level: "info",
	}
}
