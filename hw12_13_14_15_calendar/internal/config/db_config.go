package config

type DBConfig struct {
	Host            string `yaml:"host" env:"DB_HOST"`
	Port            int    `yaml:"port" env:"DB_PORT"`
	User            string `yaml:"user" env:"DB_USER"`
	Password        string `yaml:"password" env:"DB_PASSWORD"`
	DBName          string `yaml:"dbName" env:"DB_NAME"`
	MaxConn         int    `yaml:"maxConn" env:"DB_MAX_CONN" env-default:"10"`
	ItemsPerQuery   int    `yaml:"itemsPerQuery" env:"DB_ITEMS_PER_QUERY" env-default:"100"`
	RepoType        string `yaml:"repoType" env:"DB_REPO_TYPE"`
	ApplyMigrations bool   `yaml:"applyMigrations" env:"DB_APPLY_MIGRATIONS"`
	SuperUser       string `yaml:"DBSuperUser" env:"DB_SUPER_USER"`
	SuperPassword   string `yaml:"DBSuperPassword" env:"DB_SUPER_PASSWORD"`
}
