package storage

import (
	"encoding/binary"
	"fmt"
	"os"
)

type BlockReader func(*os.File, int64) error

type Block struct {
	size int32
	data []byte
	// 8 bytes
	crc []byte
	// 36 bytes
	chunkId []byte
	reader  BlockReader
}

func NewBlock(f *os.File, r BlockReader) *Block {
	// Size
	var size int32
	if err := binary.Read(f, binary.BigEndian, &size); err != nil || size <= 0 {
		return nil
	}
	// Data
	dataStart, _ := f.Seek(0, os.SEEK_CUR)
	if r != nil && r(f, int64(size)) != nil {
		fmt.Println("BlockReader Error")
		return nil
	}
	// Data & CRC
	if _, err := f.Seek(dataStart+int64(size)+CRCLen, os.SEEK_SET); err != nil {
		return nil
	}
	// ChunkId
	chunkId := make([]byte, ChunkIDLen)
	if n, err := f.Read(chunkId); err != nil || n < ChunkIDLen {
		return nil
	}
	return &Block{size, nil, nil, chunkId, r}
}

func (b Block) Size() int64 {
	return 4 + int64(b.size) + CRCLen + ChunkIDLen
}

func (b Block) String() string {
	return fmt.Sprintf("%s %d", string(b.chunkId), b.size)
}
