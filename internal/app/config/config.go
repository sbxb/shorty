package config

import "fmt"

// Config contains application settings
type Config struct {
	Proto string // should be Scheme according to RFC 3986, but not on my watch
	Host  string
	Port  string
}

// One default config to rule them all (hardcoded for now)
var DefaultConfig = Config{
	Proto: "http",
	Host:  "localhost",
	Port:  "8080",
}

// FullServerName returns "hostname:port", e.g. "localhost:8080"
func (c Config) FullServerName() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// FullServerURL returns "scheme://host:port/", e.g. "http://localhost:8080/"
func (c Config) FullServerURL() string {
	return fmt.Sprintf("%s://%s:%s/", c.Proto, c.Host, c.Port)
}
