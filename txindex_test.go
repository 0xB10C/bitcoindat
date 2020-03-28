package bitcoindat

import "testing"

func TestRead(t *testing.T) {
	b := NewBitcoinDAT(regtest0)

	tir, err := b.OpenTransactionIndexReader()
	if err != nil {
		t.Fatal(err)
	}

	ti, err := tir.ReadTXID("624d44c8b38b03fb789d9c50282bf67c5afc6959d08fc39c4f789047754a985f")
	if err != nil {
		t.Fatal(err)
	}

	if ti.NumFile != 0 || ti.DataPos != 45415 || ti.TxOffset != 1 {
		t.Errorf("Expected %+v, but got %+v", &TransactionIndex{0, 45415, 1}, ti)
	}
}
