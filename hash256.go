package bitcoindat

import (
	"encoding/hex"
)

// Hash256 represents a 32bit hash.
// These are commonly used in Bitcoin.
type Hash256 [32]byte

func (hash Hash256) String() string {
	return hex.EncodeToString(hash[:])
}

// NewHash256WithReverse takes 32bytes, reverses them and returns a Hash256.
// Bitcoin hashes are internally saved in reverse.
func NewHash256WithReverse(b [32]byte) Hash256 {
	for i := len(b)/2 - 1; i >= 0; i-- {
		opp := len(b) - 1 - i
		b[i], b[opp] = b[opp], b[i]
	}
	return Hash256(b)
}
