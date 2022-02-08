package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUserIDCookieValue(t *testing.T) {
	tests := []struct {
		cookieValue  string
		wantedResult bool
	}{
		{
			cookieValue:  "27e9ed8b869da32524db39a4bc0ea185034a33ab3f0042b5d3c81e03f6ede2c1b86ce9c1ed9b53592971729ba9312b2c",
			wantedResult: true,
		},
		{
			cookieValue:  "27e9ed8b869da32524db39a4bc0ea185034a33ab3f0042b5d3c81e03f6ede2c1b86ce9c1ed9b53592971729ba9312b2d",
			wantedResult: false,
		},
		{
			cookieValue:  "1337h4x0r",
			wantedResult: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, CheckUserIDCookieValue(tt.cookieValue), tt.wantedResult)
	}
}
