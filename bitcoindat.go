package bitcoindat

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// BitcoinDAT is a parser that can read blocks from Bitcoin Core's blk????.dat files.
type BitcoinDAT struct {
	datPath string
}

// NewBitcoinDAT is the factory function to instantiate a new BitcoinDAT
func NewBitcoinDAT(datPath string) *BitcoinDAT {
	return &BitcoinDAT{datPath: datPath}
}

func (d *BitcoinDAT) openDB(path string) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, &opt.Options{ReadOnly: true, ErrorIfMissing: true})
}
