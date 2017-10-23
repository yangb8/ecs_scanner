package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yangb8/ecsscanner/storage"
)

var (
	InvalidJLogErr = fmt.Errorf("Invalid Journal Log")
)

var (
	inputfile  string
	outputfile string
	filesize   int64
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&inputfile, "if", "", "input file")
	flag.StringVar(&outputfile, "of", "", "output file")
	flag.Int64Var(&filesize, "fs", 0, "input file size")
	flag.Parse()

	ifile, err := os.OpenFile(inputfile, os.O_RDONLY, 0444)
	checkErr(err)
	defer ifile.Close()

	size := filesize
	if size == 0 {
		info, err := ifile.Stat()
		checkErr(err)
		size = info.Size()
	}

	ofile, err := os.OpenFile(outputfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	defer ofile.Close()
	writter := bufio.NewWriter(ofile)
	defer writter.Flush()
	log.SetOutput(writter)

	var offset int64
	for offset = 0; offset+storage.ChunkSize <= size; offset += storage.ChunkSize {
		fmt.Println("##", inputfile, offset, size)
		ifile.Seek(offset, os.SEEK_SET)
		storage.NewChunk(ifile, offset, storage.ChunkSize, chunkScanner).Scan()
	}
}

func chunkScanner(f *os.File, offset, length int64, id []byte) error {
	n, _ := f.Seek(0, os.SEEK_CUR)
	log.Printf("%s %s %d %d %d", string(id), inputfile, n/storage.ChunkSize*storage.ChunkSize, n%storage.ChunkSize, length)
	return nil
}
