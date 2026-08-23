package storage

type Config struct {
	Port           int
	SqliteDSN      string
	FileStorageDir string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		Port:           8001,
		SqliteDSN:      "storage.sqlite3",
		FileStorageDir: ".storage",
	}
	return cfg, nil
}
