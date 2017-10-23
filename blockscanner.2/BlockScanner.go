package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/yangb8/ecsscanner/or"
	"github.com/yangb8/ecsscanner/storage"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	InvalidJLogErr = fmt.Errorf("Invalid Journal Log")
)

var (
	inputfile    string
	outputfile   string
	threadnumber int
	re           *regexp.Regexp = regexp.MustCompile("[a-f0-9]{64}")
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&inputfile, "if", "", "input location file")
	flag.StringVar(&outputfile, "of", "", "output file")
	flag.IntVar(&threadnumber, "n", 1, "thread number")
	flag.Parse()

	ifile, err := os.OpenFile(inputfile, os.O_RDONLY, 0444)
	checkErr(err)
	defer ifile.Close()

	log.SetOutput(&lumberjack.Logger{
		Filename: outputfile,
		MaxSize:  512, // megabytes
	})

	c := producer(ifile)

	var wg sync.WaitGroup
	for i := 0; i < threadnumber; i++ {
		wg.Add(1)
		go func(idx int, ch <-chan *entry) {
			defer wg.Done()
			fmt.Println("goroutine ", idx, " started")
			for e := range ch {
				processor(e)
			}
			fmt.Println("goroutine ", idx, " stopped")
		}(i, c)
	}
	wg.Wait()
}

func producer(f *os.File) <-chan *entry {
	c := make(chan *entry)

	r := bufio.NewReader(f)
	go func() {
		defer close(c)
		for i := 1; ; i++ {
			s, err := r.ReadString('\n')
			if err != nil {
				break
			}
			if parts := strings.Split(s, " "); len(parts) == 2 {
				f := strings.TrimSpace(parts[0])
				n, _ := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
				c <- &entry{i, f, n}
			}
		}
	}()

	return c
}

type entry struct {
	seq   int
	fname string
	size  int64
}

func processor(e *entry) {
	ifile, err := os.OpenFile(e.fname, os.O_RDONLY, 0444)
	if err != nil {
		return
	}
	defer ifile.Close()

	var offset int64
	for offset = 0; offset+storage.ChunkSize <= e.size; offset += storage.ChunkSize {
		fmt.Printf("## No.%d %s: %d/%d\n", e.seq, e.fname, offset/storage.ChunkSize, e.size/storage.ChunkSize)
		ifile.Seek(offset, os.SEEK_SET)
		storage.NewChunk(ifile, offset, storage.ChunkSize, chunkScanner).Scan()
	}
}

func chunkScanner(f *os.File, offset, length int64, id []byte) error {
	var b *storage.Block
	for p := offset; p < length; {
		if b = storage.NewBlock(f, JournalLogReader); b == nil {
			break
		}
		p += b.Size()
	}
	return nil
}

func JournalLogReader(f *os.File, length int64) error {
	var (
		offset int64
		size   int32
	)

	startOffset, _ := f.Seek(0, os.SEEK_CUR)
	buffer := make([]byte, 2*storage.MB)

	for offset < length {
		f.Seek(startOffset+offset, os.SEEK_SET)
		if offset += 4; offset > length {
			return InvalidJLogErr
		}
		if err := binary.Read(f, binary.BigEndian, &size); err != nil || size <= 0 || size > 2*storage.MB {
			return InvalidJLogErr
		}

		if offset += int64(size); offset > length {
			return InvalidJLogErr
		}
		if n, err := f.Read(buffer[:size]); err != nil || n < int(size) {
			return InvalidJLogErr
		}

		header := &or.JournalDirTableLogHeader{}
		if err := proto.Unmarshal(buffer[:size], header); err != nil || header.GetPayloadLength() < 0 {
			fmt.Println("Decode error", err)
			return InvalidJLogErr
		}

		if offset += int64(header.GetPayloadLength()); offset > length {
			return InvalidJLogErr
		}

		if header.GetPayloadLength() > 0 && header.GetPayloadLength() <= 2*storage.MB && header.GetSchemaKey().GetType() == or.SchemaKeyType_OBJECT_TABLE_KEY && header.GetSchemaKey().GetUserKey() != nil {
			key := re.FindString(string(header.GetSchemaKey().GetUserKey()))
			if key == "" {
				continue
			}

			if n, err := f.Read(buffer[:header.GetPayloadLength()]); err != nil || n < int(header.GetPayloadLength()) {
				continue
			}

			rec := &or.DirectoryUpdateRecord{}
			if err := proto.Unmarshal(buffer[:header.GetPayloadLength()], rec); err == nil {
				var (
					locs    []*or.SegmentLocation
					indices []*or.DataIndex
				)

				if rec.GetSegment().GetSegmentUMR() != nil {
					locs = rec.GetSegment().GetSegmentUMR().GetReposUMRLocations()
					indices = rec.GetSegment().GetSegmentUMR().GetDataIndices()
				} else if rec.GetSegment().GetSegmentIMR() != nil {
					locs = rec.GetSegment().GetSegmentIMR().GetReposIMRLocations()
					indices = rec.GetSegment().GetSegmentIMR().GetDataIndices()
				}

				var msg string
				if locs != nil {
					msg += fmt.Sprintf("NLOC %d ", len(locs))
					for i, loc := range locs {
						msg += fmt.Sprintf("loc_%d %s %d %d %d %d ", i, loc.GetChunkId(), loc.GetOffset(), loc.GetEndOffset(), loc.GetRangeInfo().GetRelativeOffset(), loc.GetRangeInfo().GetRelativeEndOffset())
					}
				}
				if indices != nil {
					msg += fmt.Sprintf("NIDX %d ", len(indices))
					for i, index := range indices {
						msg += fmt.Sprintf("idx_%d %d %d ", i, index.GetDataRange().GetObjectOffset(), index.GetDataRange().GetObjectLength())
					}
				}
				if msg != "" {
					log.Printf("%s - %s ", key, msg)
				}
			}
		}
	}
	return nil
}
