package bitcoindat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sort"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	// BlockValidHeader - Parsed, version ok, hash satisfies claimed PoW, 1 <= vtx count <= max, timestamp not in future
	BlockValidHeader = 1
	// BlockValidTree - All parent headers found, difficulty matches, timestamp >= median previous, checkpoint. Implies all parents are also at least TREE.
	BlockValidTree = 2
	// BlockValidTransactions - Only first tx is coinbase, 2 <= coinbase input script length <= 100, transactions valid, no duplicate txids, sigops, size, merkle root. Implies all parents are at least TREE but not necessarily TRANSACTIONS. When all parent blocks also have TRANSACTIONS.
	BlockValidTransactions = 3
	// BlockValidChain - Outputs do not overspend inputs, no double spends, coinbase output ok, immature coinbase spends, BIP30. Implies all parents are also at least CHAIN.
	BlockValidChain = 4
	// BlockValidScripts - Scripts & signatures ok. Implies all parents are also at least SCRIPTS.
	BlockValidScripts = 5
	// BlockValidMask - All validity bits.
	BlockValidMask = BlockValidHeader | BlockValidTree | BlockValidTransactions | BlockValidChain | BlockValidScripts
	// BlockHaveData - full block available in blk*.dat
	BlockHaveData = 8
	// BlockHaveUndo - undo data available in rev*.dat
	BlockHaveUndo = 16
	// BlockHaveMask - Block and Undo data avaliable
	BlockHaveMask = BlockHaveData | BlockHaveUndo
	// BlockFailedValid - stage after last reached validness failed
	BlockFailedValid = 32
	// BlockFailedChild - descends from failed block
	BlockFailedChild = 64
	// BlockFailedMask - Mask for failed block.
	BlockFailedMask = BlockFailedValid | BlockFailedChild
)

// BlockIndices is a list of multiple BlockIndex which can contain multiple blocks at the same hight and forks
type BlockIndices []BlockIndex

// BlockIndexChain is a list of multiple BlockIndex which has only one block per hight
type BlockIndexChain []BlockIndex

// BlockIndex represents an index entry in the leveldb database.
type BlockIndex struct {
	Hash               Hash256 // Block Hash
	BitcoinCoreVersion int32   // Bitcoin Core Version
	Height             int32   // Block height
	Status             uint32  // Block status
	NumTx              uint32  // Number of transactions
	NumFile            int32   // File Number
	PosData            uint32  // Position of the block data in the file
	PosUndo            uint32  // Position of the undo data for that block in the undo file
	BlockVersion       int32   // Block version
	PreviousHash       Hash256 // Previous Block Hash
	MerkleRootHash     Hash256 // Hash of the Merkle Root
	Time               uint32  // Block time
	Bits               uint32  // Bits
	Nonce              uint32  // Nonce
}

// GetBlockIndices reads all indexed blocks from the leveldb database and
// returns these unordered.
func (p *BitcoinDAT) GetBlockIndices() (BlockIndices, error) {
	db, err := p.openDB(p.datadir + "/blocks/index")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	bi, err := buildBlockIndex(db)
	if err != nil {
		return nil, err
	}

	return bi, nil
}

func buildBlockIndex(db *leveldb.DB) (BlockIndices, error) {
	bi := make([]BlockIndex, 0)
	iter := db.NewIterator(util.BytesPrefix([]byte("b")), nil)
	for iter.Next() {
		ib, err := parseBlockIndex(iter.Key(), iter.Value())
		if err != nil {
			return nil, err
		}
		bi = append(bi, ib)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return bi, nil
}

// GetMainChain returns the most work chain of blocks containing no forks.
// BlockIndexChain is sorted by ascending block height (genesis -> tip).
func (bi *BlockIndices) GetMainChain() BlockIndexChain {
	blocksAtHeight := make(map[int]BlockIndices)
	for _, b := range *bi {
		h := int(b.Height)
		if blocksAtHeight[h] == nil {
			blocksAtHeight[h] = make(BlockIndices, 0)
		}
		blocksAtHeight[h] = append(blocksAtHeight[h], b)
	}

	// create a list of heights starting that the tip and going
	// down to genesis
	keys := make([]int, 0, len(blocksAtHeight))
	for k := range blocksAtHeight {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// create a list of block indexes representing the main or most-work
	// chain.
	// the list is created in reverse. starting with the tip and then adding
	// the block at the next lower hight. if there are multiple blocks at a
	// lower hight then pick the one that the previous one refereed to.
	bic := BlockIndexChain{}
	for _, height := range keys {
		blocksAtCurHeight := blocksAtHeight[height]
		if len(blocksAtCurHeight) == 1 {
			bic = append(bic, blocksAtCurHeight[0])
		} else {
			// if there are two or more blocks for a given height we choose
			// the block the previous block refereed to
			if len(bic) > 0 {
				didAppend := false
				for _, b := range blocksAtCurHeight {
					if bytes.Equal(b.Hash[:], bic[len(bic)-1].PreviousHash[:]) {
						bic = append(bic, b)
						didAppend = true
					}
				}
				// check that we had something to append. this only ever should be
				// true if there is an error in the database or our deserialization
				// of the data block index in the database.
				if !didAppend {
					panic(fmt.Sprintf("No previous Block found with the hash %s at height %d!", bic[len(bic)-1].PreviousHash, height))
				}
			} else {
				// just choose the first block we have at the highest height as tip
				// FIXME: use chainwork
				bic = append(bic, blocksAtCurHeight[0])
			}
		}
	}
	// reverse bic to start at genesis -> tip
	for i := len(bic)/2 - 1; i >= 0; i-- {
		opp := len(bic) - 1 - i
		bic[i], bic[opp] = bic[opp], bic[i]
	}

	return bic
}

func parseBlockIndex(key []byte, value []byte) (ib BlockIndex, err error) {
	hashBytes := [32]byte{}
	copy(hashBytes[:], key[1:33])
	hash := NewHash256(hashBytes).ReversedCopy()

	data := NewMovingBytes(value)

	ib = BlockIndex{}
	ib.Hash = hash
	ib.BitcoinCoreVersion = int32(data.ReadVarInt())
	ib.Height = int32(data.ReadVarInt())
	ib.Status = uint32(data.ReadVarInt())
	ib.NumTx = uint32(data.ReadVarInt())

	if ib.Status&BlockHaveMask > 0 {
		ib.NumFile = int32(data.ReadVarInt())
	}

	if ib.Status&BlockHaveData > 0 {
		ib.PosData = uint32(data.ReadVarInt())
	}

	if ib.Status&BlockHaveUndo > 0 {
		ib.PosUndo = uint32(data.ReadVarInt())
	}

	ib.BlockVersion = data.ReadInt32LitteEndian()
	ib.PreviousHash = data.ReadHash()
	ib.MerkleRootHash = data.ReadHash()
	ib.Time = data.ReadUInt32LitteEndian()
	ib.Bits = data.ReadUInt32LitteEndian()
	ib.Nonce = data.ReadUInt32LitteEndian()

	return
}

// ReadBlockData reads the data for the passed BlockIndex and returns a byte slice
// with the raw block data.
func (p *BitcoinDAT) ReadBlockData(ib BlockIndex) ([]byte, error) {
	if ib.Status&BlockHaveData == 0 {
		return nil, fmt.Errorf("No data avaliable for this block")
	}

	filename := fmt.Sprintf("blk%05d.dat", ib.NumFile)
	file, err := os.Open(p.datadir + "blocks/" + filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// read block size
	sizeBytes := make([]byte, 4)
	_, err = file.Seek(int64(ib.PosData)-4, 0)
	if err != nil {
		return nil, fmt.Errorf("could not seek the block size: %w", err)
	}
	_, err = file.Read(sizeBytes)
	if err != nil {
		return nil, fmt.Errorf("could not read the block size: %w", err)
	}

	size := binary.LittleEndian.Uint32(sizeBytes)

	block := make([]byte, size)
	_, err = file.Seek(int64(ib.PosData), 0)
	if err != nil {
		return nil, fmt.Errorf("could not seek the block position: %w", err)
	}
	_, err = file.Read(block)
	if err != nil {
		return nil, fmt.Errorf("could not read the block: %w", err)
	}

	return block, nil
}
