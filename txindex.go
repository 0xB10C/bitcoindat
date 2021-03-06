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
	NumFile  int32
	DataPos  uint32
	TxOffset int
}

// OpenTransactionIndexReader opens a transaction index level db.
// The DB must be closed after use, by calling Close method.
func (p *BitcoinDAT) OpenTransactionIndexReader() (*TransactionIndexReader, error) {
	db, err := p.openDB(p.datadir + "/indexes/txindex")
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

	if err != nil {
		return nil, err
	}

	return tir.Read(NewHash256ByteSlice(b))
}

// Read reads a TransactionIndex from the LevelDB
func (tir *TransactionIndexReader) Read(txid Hash256) (*TransactionIndex, error) {

	// Bitcoin Core handles and saves the transaction ids in reverse.
	txidReverse := txid.ReversedCopy()

	key := append([]byte("t"), txidReverse[:]...)
	value, err := tir.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	mb := NewMovingBytes(value)
	ti := &TransactionIndex{}
	ti.NumFile = int32(mb.ReadVarInt())
	ti.DataPos = uint32(mb.ReadVarInt())
	ti.TxOffset = int(mb.ReadVarInt())

	return ti, nil
}
