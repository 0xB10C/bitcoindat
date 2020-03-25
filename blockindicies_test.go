package bitcoindat

import "testing"

const regtest0 = "./testdata/regtest0"

func TestBlockIndex(t *testing.T) {
	b := NewBitcoinDAT(regtest0)

	bi, err := b.GetBlockIndices()
	if err != nil {
		t.Fatal(err)
	}

	const numBlockIndexes = 206
	if len(bi) != numBlockIndexes {
		t.Errorf("The %s block indexes should have %d entries; got %d", regtest0, numBlockIndexes, len(bi))
	}

	mc := bi.GetMainChain()

	const numBlocksMainChain = 201
	if len(mc) != numBlocksMainChain {
		t.Errorf("The %s main chain should have %d entries; got %d", regtest0, numBlocksMainChain, len(mc))
	}
}
