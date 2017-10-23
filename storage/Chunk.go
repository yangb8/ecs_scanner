package storage

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type ChunkHandler func(*os.File, int64, int64, []byte) error

type Chunk struct {
	f           *os.File
	startOffset int64
	length      int64
	handler     ChunkHandler
}

func NewChunk(f *os.File, startOffset, length int64, hlr ChunkHandler) *Chunk {
	return &Chunk{f, startOffset, length, hlr}
}

func (c Chunk) getFirstChunkId() ([]byte, error) {
	if _, err := c.f.Seek(c.startOffset, os.SEEK_SET); err != nil {
		return nil, err
	}
	// try first block
	block := NewBlock(c.f, nil)
	if block == nil {
		return nil, InvalidChunkErr
	}
	if isEmptyChunkId(block.chunkId) {
		return nil, EmptyChunkIDErr
	}

	return block.chunkId, nil
}

func (c Chunk) getLastChunkId() ([]byte, error) {
	if c.length < ChunkIDLen {
		return nil, InvalidChunkErr
	}

	buffer := make([]byte, 1*MB)
	var offset, pos int64
	for offset, pos = c.startOffset+ChunkSize-1*MB, c.startOffset+ChunkSize-1; offset >= 0; offset -= 1 * MB {
		if n, err := c.f.ReadAt(buffer, offset); err != nil || n < 1*MB {
			return nil, InvalidChunkErr
		}
		var i int
		for i = 1*MB - 1; i >= 0 && buffer[i] == 0; i-- {
			pos--
		}
		if i >= 0 {
			break
		}
	}

	if pos < c.startOffset+ChunkIDLen-1 {
		return nil, InvalidChunkIDErr
	}

	chunkId := make([]byte, ChunkIDLen)
	if n, err := c.f.ReadAt(chunkId, pos-ChunkIDLen+1); err != nil || n < ChunkIDLen {
		return nil, InvalidChunkIDErr
	}

	return chunkId, nil
}

func (c Chunk) findChunkId(id []byte) (int64, error) {
	if id == nil || len(id) != ChunkIDLen {
		return -1, InvalidChunkIDErr
	}

	buffer := make([]byte, 1*MB)
	var offset, pos int64
	for offset, pos = c.startOffset, 0; offset < c.startOffset+c.length; offset, pos = offset+1*MB-36, pos+1*MB-36 {
		if offset+1*MB > c.startOffset+c.length {
			buffer = buffer[:c.startOffset+c.length-offset]
		}
		if n, err := c.f.ReadAt(buffer, offset); err != nil || n < len(buffer) {
			return -1, InvalidChunkErr
		}
		if idx := strings.Index(string(buffer), string(id)); idx >= 0 {
			return pos + int64(idx), nil
		}
	}
	return -1, fmt.Errorf("Not Found")
}

func (c *Chunk) getValidStart() (int64, []byte, error) {
	var (
		firstID, lastID []byte
		pos             int64
		err             error
	)
	if firstID, err = c.getFirstChunkId(); err != nil {
		return -1, nil, err
	}
	if lastID, err = c.getLastChunkId(); err != nil {
		return -1, nil, err
	}

	if reflect.DeepEqual(firstID, lastID) {
		return 0, lastID, nil
	}

	if pos, err = c.findChunkId(lastID); err != nil {
		return -1, nil, err
	}
	return pos + ChunkIDLen, lastID, nil
}

func (c Chunk) Scan() (err error) {
	var offset int64
	var id []byte
	if offset, id, err = c.getValidStart(); err != nil {
		return err
	}

	c.f.Seek(c.startOffset+offset, os.SEEK_SET)

	if c.handler != nil {
		c.handler(c.f, offset, c.length, id)
	}
	return nil
}

func isEmptyChunkId(id []byte) bool {
	for _, b := range id {
		if b != 0 {
			return false
		}
	}
	return true
}
