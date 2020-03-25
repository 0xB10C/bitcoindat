package datparser

import (
	"encoding/binary"
	"fmt"
	"os"
)

// DATParser is a parser that can read blocks from Bitcoin Core's blk????.dat files.
type DATParser struct {
	datPath string
}

// NewDATParser is the factory function to instantiate a new DATParser
func NewDATParser(datPath string) (*DATParser, error) {
	return &DATParser{datPath: datPath}, nil
}

// ReadBlockData reads the data for the passed BlockIndex and returns a byte slice
// with the raw block data.
func (d *DATParser) ReadBlockData(ib BlockIndex) ([]byte, error) {
	if ib.Status&BlockHaveData == 0 {
		return nil, fmt.Errorf("No data avaliable for this block")
	}

	filename := fmt.Sprintf("blk%05d.dat", ib.NumFile)
	file, err := os.Open(d.datPath + "/" + filename)
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
