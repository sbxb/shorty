package url

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"strings"
)

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	Result string `json:"result"`
}

type URLEntry struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchURLRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

//type BatchURLResponse []BatchURLResponseEntry

type BatchURLEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"-"`
	ShortURL      string `json:"short_url"`
}

//type DeleteList []string

// ShortID converts URL to a string containing its MD5 hash represented
// as a base62 number
// MD5 is OK since we do not care about any (almost impossible) collisions
// NB: Go base62 implementation (using digits 0..9a..zA..Z) may differ from any
// other implementations (using digits 0..9A..Za..z)
func ShortID(url string) string {
	hash := md5.Sum([]byte(url))
	strHash := hex.EncodeToString(hash[:])

	n := new(big.Int)
	// strHash always contains a valid 128-bit hexadecimal number (as a string)
	n.SetString(strHash, 16)

	return n.Text(62)
}

func IsBatchURLRequestValid(batch []BatchURLRequestEntry) bool {
	if len(batch) == 0 {
		return false
	}
	for _, u := range batch {
		if u.CorrelationID == "" || !IsValidInputURL(u.OriginalURL) {
			return false
		}
	}

	return true
}

// IsValidInputURL checks if the user input slightly resembles a valid URL or not
// by simply detecting non-valid characters
// There is no need to parse a URL, let the user shorten whatever they want
// within reasonable limits
func IsValidInputURL(url string) bool {
	// characters that URL can possibly contain according to RFC3986
	allowedChars := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789:/?#[]@!$&'()*+,;=-_.~%`

	// empty or extremely long - not valid
	if url == "" || len(url) > 2048 {
		return false
	}

	for _, c := range url {
		if !strings.Contains(allowedChars, string(c)) {
			return false
		}
	}
	return true
}
