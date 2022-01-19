package config

import (
	"flag"
	"os"
	"strings"
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

// New creates config by merging default settings with flags, then with env variables
// The last nonempty value takes precedence (default < flag < env) except for
// FILE_STORAGE_PATH env variable which overrides -f flag even if empty
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

	fsp, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		// empty string is valid here, overrides -f flag and returns the default ""
		c.FileStoragePath = fsp
	}
}

func (c *Config) Validate() error {
	// Remove leading and trailing spaces without complaining
	// Other mistakes and typos are to be considered as errors
	c.ServerAddress = strings.TrimSpace(c.ServerAddress)
	c.BaseURL = strings.TrimSpace(c.BaseURL)
	c.FileStoragePath = strings.TrimSpace(c.FileStoragePath)

	if err := ValidateServerAddress(c.ServerAddress); err != nil {
		return err
	}

	if err := ValidateBaseURL(c.BaseURL); err != nil {
		return err
	}

	// No need to validate c.FileStoragePath, storage itself will do the job
	return nil
}
