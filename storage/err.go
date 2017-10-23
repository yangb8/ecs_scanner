package storage

import "fmt"

var (
	InvalidBlockErr   = fmt.Errorf("Invalid Block ID")
	InvalidChunkErr   = fmt.Errorf("Invalid Chunk")
	InvalidChunkIDErr = fmt.Errorf("Invalid Chunk ID")
	EmptyChunkIDErr   = fmt.Errorf("Empty Chunk ID")
)
