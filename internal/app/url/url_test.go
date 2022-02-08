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

func TestValidateInputURL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "http://www.example.com/example",
			want:  true,
		},
		{
			input: "http://www.example.com/example/?id=15&uid=25#top",
			want:  true,
		},
		{
			input: "https://yandex.ru/pogoda/?utm_source=home&utm_campaign=informer&utm_medium=web&utm_content=main_informer&utm_term=current_day_part",
			want:  true,
		},
		{
			input: "https://cs.opensource.google/go/go/+/refs/tags/go1.17.6:src/strings/strings.go;l=66",
			want:  true,
		},
		{
			input: "       ",
			want:  false,
		},
		{
			input: "яндекс.рф",
			want:  false,
		},
		{
			input: "Once upon a time there lived ...",
			want:  false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, url.IsValidInputURL(tt.input), tt.want)
	}
}
