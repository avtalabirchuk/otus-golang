package config

type LogConfig struct {
	Level    string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	FilePath string `yaml:"filePath" env:"LOG_FILE"`
}
