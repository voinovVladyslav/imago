package gateway

type Config struct {
	Port        int
	FileRepoURL string
}

func NewConfig() (*Config, error) {
	cfg := &Config{Port: 8000, FileRepoURL: "http://localhost:8001"}
	return cfg, nil
}
