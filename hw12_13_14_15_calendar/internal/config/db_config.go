package config

type DBConfig struct {
	Host          string `yaml:"host" config:"required"`
	Port          int    `yaml:"port" config:"required"`
	DBName        string `yaml:"dbName" config:"required"`
	User          string `yaml:"user" config:"required"`
	Pass          string `yaml:"pass" config:"required"`
	MaxConn       int    `yaml:"maxConn" config:"required"`
	ItemsPerQuery int    `yaml:"itemsPerQuery" config:"required"`
	RepoType      string `yaml:"repoType" config:"required"`
}

func defaultDBConfig() DBConfig {
	return DBConfig{
		Host:          "localhost",
		Port:          5432,
		MaxConn:       10,
		ItemsPerQuery: 100,
	}
}
