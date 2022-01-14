package config

import (
	"flag"
	"os"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
)

// Config contains application settings
type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

var defaultConfig = Config{
	ServerAddress: defaultServerAddress,
	BaseURL:       defaultBaseURL,
}

// New creates config merging default settings with flags, then with env variables
// The last nonempty value takes precedence (default < flag < env)
func New() *Config {
	c := defaultConfig
	c.parseFlags()
	c.parseEnvVars()

	return &c
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", defaultServerAddress, "network address the server listens on")
	flag.StringVar(&c.BaseURL, "b", defaultBaseURL, "resulting base URL")
	flag.StringVar(&c.FileStoragePath, "f", "", `storage file (default "")`)

	flag.Parse()
}

func (c *Config) parseEnvVars() {
	sa := os.Getenv("SERVER_ADDRESS")
	if sa != "" {
		c.ServerAddress = sa
	}

	bu := os.Getenv("BASE_URL")
	if bu != "" {
		c.BaseURL = bu
	}

	fsp := os.Getenv("FILE_STORAGE_PATH")
	if fsp != "" {
		c.FileStoragePath = fsp
	}
}
