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
	DatabaseDSN     string
}

var defaultConfig = Config{
	ServerAddress: defaultServerAddress,
	BaseURL:       defaultBaseURL,
}

// New creates config by merging default settings with flags, then with env variables
// The last nonempty value takes precedence (default < flag < env) except for
// FILE_STORAGE_PATH / DATABASE_DSN env variables which overrides -f / -d flags
// even if empty
// New also handles validation and returns non-nil error if validation failed
func New() (Config, error) {
	c := defaultConfig
	c.parseFlags()
	c.parseEnvVars()
	err := c.Validate()
	return c, err
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", defaultServerAddress, "network address the server listens on")
	flag.StringVar(&c.BaseURL, "b", defaultBaseURL, "resulting base URL")
	flag.StringVar(&c.FileStoragePath, "f", "", `storage file (default "")`)
	flag.StringVar(&c.DatabaseDSN, "d", "", `database dsn (default "")`)

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

	dd, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		// empty string is valid here, overrides -d flag and returns the default ""
		c.DatabaseDSN = dd
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
