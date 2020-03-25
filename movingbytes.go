package bitcoindat

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type MovingBytes struct {
	b   []byte
	pos uint64
}

func NewMovingBytes(b []byte) MovingBytes {
	return MovingBytes{b, 0}
}

func (mb *MovingBytes) ReadVarInt() (n uint64) {
	for true {
		b := mb.b[mb.pos : mb.pos+1][0]
		mb.pos++
		n = (n << uint64(7)) | uint64(b&uint8(0x7F))
		if b&uint8(0x80) > 0 {
			n++
		} else {
			return
		}
	}
	return
}

func (mb *MovingBytes) ReadUInt32LitteEndian() (val uint32) {
	val = binary.BigEndian.Uint32(mb.b[mb.pos : mb.pos+4])
	mb.pos += 4
	return
}

func (mb *MovingBytes) PrintNext(n uint64, tag string) {
	fmt.Printf("Next at pos %d (%s): %s\n", mb.pos, tag, hex.EncodeToString(mb.b[mb.pos:mb.pos+n]))
}

func (mb *MovingBytes) ReadInt32LitteEndian() (val int32) {
	val = int32(binary.BigEndian.Uint32(mb.b[mb.pos : mb.pos+4]))
	mb.pos += 4
	return
}

func (mb *MovingBytes) ReadHash() (hash Hash256) {
	copy(hash[:], mb.b[mb.pos:mb.pos+32])
	// reverse
	for i := len(hash)/2 - 1; i >= 0; i-- {
		opp := len(hash) - 1 - i
		hash[i], hash[opp] = hash[opp], hash[i]
	}
	mb.pos += 32
	return
}
