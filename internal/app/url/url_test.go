package url_test

import (
	"testing"

	"github.com/sbxb/shorty/internal/app/url"
)

func TestShortId(t *testing.T) {
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
		id := url.ShortId(tt.original)
		if id != tt.want {
			t.Errorf("got [%s], want [%s]", id, tt.want)
		}
	}

}
