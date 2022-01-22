package url_test

import (
	"testing"

	"github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
)

func TestShortID(t *testing.T) {
	tests := []struct {
		original string
		want     string
	}{
		{
			original: "http://www.example.com/example",
			want:     "3r6MoXVrMDuxcVr2I70g30",
		},
		{
			original: "jk8ssl", // md5 sum with leading zeroes "00000000 18e6137a c2caab16 074784a6"
			want:     "a1lqzNaaPw3Kxzpk",
		},
	}

	for _, tt := range tests {
		id := url.ShortID(tt.original)

		assert.Equal(t, id, tt.want)
	}
}
