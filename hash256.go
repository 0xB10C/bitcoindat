package bitcoindat

import (
	"encoding/hex"
	"fmt"
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

// NewHash256 takes 32bytes returns a Hash256.
func NewHash256(b [32]byte) Hash256 {
	return Hash256(b)
}

// NewHash256FromByteSlice creates a Hash256 from a bytes slice. An error is
// returned if the byte slice is not 32 bytes long.
func NewHash256FromByteSlice(b []byte) (Hash256, error) {
	if len(b) != 32 {
		return Hash256{}, fmt.Errorf("hash256: need 32 bytes, got %d bytes", len(b))
	}

	var hash [32]byte
	copy(hash[:], b[:32])
	return Hash256(hash), nil
}
