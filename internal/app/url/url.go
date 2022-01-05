package url

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
)

// ShortId converts URL to a string containing its MD5 hash represented
// as a base62 number
// MD5 is OK since we do not care about any (almost impossible) collisions
// NB: Go base62 implementation (using digits 0..9a..zA..Z) may differ from any
// other implementations (using digits 0..9A..Za..z)
func ShortId(url string) string {
	hash := md5.Sum([]byte(url))
	strHash := hex.EncodeToString(hash[:])

	n := new(big.Int)
	// strHash always contains a valid 128-bit hexadecimal number (as a string)
	n.SetString(strHash, 16)

	return n.Text(62)
}
