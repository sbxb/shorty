package config

import (
	"fmt"
	"os"
)

// Config contains application settings
type Config struct {
	proto           string // should be Scheme according to RFC 3986, but not on my watch
	host            string
	port            string
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

// One default config to rule them all (hardcoded for now)
var defaultConfig = Config{
	proto: "http",
	host:  "localhost",
	port:  "8080",
}

// defaultServerName returns "hostname:port", e.g. "localhost:8080"
func (c Config) defaultServerName() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

// defaultServerURL returns "scheme://host:port", e.g. "http://localhost:8080"
func (c Config) defaultServerURL() string {
	return fmt.Sprintf("%s://%s:%s", c.proto, c.host, c.port)
}

// New creates new config merging env variables with the default ones
func New() *Config {
	c := defaultConfig

	sa := os.Getenv("SERVER_ADDRESS")
	if sa == "" {
		sa = c.defaultServerName()
	}
	c.ServerAddress = sa

	bu := os.Getenv("BASE_URL")
	if bu == "" {
		bu = c.defaultServerURL()
	}
	c.BaseURL = bu

	c.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	return &c
}
