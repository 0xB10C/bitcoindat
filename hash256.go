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

// ReversedCopy returns a reversed copy of the Hash256
// Bitcoin Core handles hashes internally and in the LevelDB in reverse.
func (hash Hash256) ReversedCopy() Hash256 {
	var reversed [32]byte
	copy(reversed[:], hash[:32])
	for i := len(reversed)/2 - 1; i >= 0; i-- {
		opp := len(reversed) - 1 - i
		reversed[i], reversed[opp] = reversed[opp], reversed[i]
	}
	return reversed
}

// NewHash256 takes 32bytes returns a Hash256.
func NewHash256(b [32]byte) Hash256 {
	return Hash256(b)
}

// NewHash256ByteSlice creates a Hash256 from a bytes slice.
func NewHash256ByteSlice(b []byte) Hash256 {
	if len(b) != 32 {
		panic("hash256: expected byte slice of length 32")
	}

	var hash [32]byte
	copy(hash[:], b[:32])
	return NewHash256(hash)
}
