package url

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
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
		if u.CorrelationID == "" || u.OriginalURL == "" {
			return false
		}
	}

	return true
}
