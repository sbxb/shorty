package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBaseURL_Valid(t *testing.T) {
	tests := []string{
		"http://localhost:8080",
		"http://127.0.0.1:5555",
		"https://localhost",
		"https://127.0.0.1",
		"ftp://example.com/public",
		"http://v3ry.l0ng.example.com:443",
		"https://xn--80atccmdviy.xn--e1aybc",
	}

	for _, tt := range tests {
		err := ValidateBaseURL(tt)
		assert.Nil(t, err)
	}
}

func TestValidateBaseURL_NotValid(t *testing.T) {
	tests := []string{
		"",
		":8080",
		"localhost:8080",
		"127.0.0.1.5555",
		"https://",
		"https://280.0.0.1",
		"https://280.1",
		"ftp://example..com..public",
		"ftp:/example.com/public",
		"http://v3ry.l0ng.example.com:443210",
		"https://--80atccmdviy.--e1aybc",
	}

	for _, tt := range tests {
		err := ValidateBaseURL(tt)
		assert.NotNil(t, err)
	}
}

func TestValidateServerAddress_Valid(t *testing.T) {
	tests := []string{
		"localhost:8080",
		"127.0.0.1:5555",
		":80",
		"example.com:443",
		"v3ry.l0ng.example.com:443",
		"xn--80atccmdviy.xn--e1aybc:5555",
	}

	for _, tt := range tests {
		err := ValidateServerAddress(tt)
		assert.Nil(t, err)
	}

}

func TestValidateServerAddress_NotValid(t *testing.T) {
	tests := []string{
		"",
		" ",
		":",
		":801234",
		"localhost",
		"www.example.com",
		"280.1.1.1:21",
		"280.280.280.280.280:280",
	}

	for _, tt := range tests {
		err := ValidateServerAddress(tt)
		assert.NotNil(t, err)
	}
}
