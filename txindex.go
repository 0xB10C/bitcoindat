package bitcoindat

import (
	"encoding/hex"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// TransactionIndexReader reads TransactionIndex's from a LevelDB
type TransactionIndexReader struct {
	db *leveldb.DB
}

// TransactionIndex holds the block position of a transaction
type TransactionIndex struct {
	blockFile int
	blockPos  int
	txOffset  int
}

// OpenTransactionIndexReader opens a transaction index level db.
// The DB must be closed after use, by calling Close method.
func (p *BitcoinDAT) OpenTransactionIndexReader() (*TransactionIndexReader, error) {
	db, err := p.openDB(p.datPath + "/indexes/txindex")
	if err != nil {
		return nil, err
	}

	return &TransactionIndexReader{db}, nil
}

// Close closes the underlying LevelDB
func (tir *TransactionIndexReader) Close() {
	tir.db.Close()
}

// ReadTXID reads a TransactionIndex from the LevelDB. This function wraps Read()
func (tir *TransactionIndexReader) ReadTXID(txid string) (*TransactionIndex, error) {
	if len(txid) != 64 {
		return nil, fmt.Errorf("Expected txid string to have a length of 64 chars; got %d for txid: %s", len(txid), txid)
	}

	b, err := hex.DecodeString(txid)
	if err != nil {
		return nil, err
	}

	return tir.Read(b)
}

// Read reads a TransactionIndex from the LevelDB
func (tir *TransactionIndexReader) Read(txid []byte) (*TransactionIndex, error) {
	if len(txid) != 32 {
		return nil, fmt.Errorf("Expected txid byte array to have a length of 32 byte; got %d for txid: %s", len(txid), hex.EncodeToString(txid))
	}

	var tmp [32]byte
	copy(tmp[:], txid)
	txid256 := NewHash256WithReverse(tmp)

	key := append([]byte("t"), txid256[:]...)
	value, err := tir.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	mb := NewMovingBytes(value)
	ti := &TransactionIndex{}
	ti.blockFile = int(mb.ReadVarInt())
	ti.blockPos = int(mb.ReadVarInt())
	ti.txOffset = int(mb.ReadVarInt())

	return ti, nil
}
