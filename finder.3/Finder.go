package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/yangb8/ecsscanner/storage"
)

var (
	obf    string
	c2f    string
	o2c    string
	outdir string
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&obf, "obf", "", "objects file")
	flag.StringVar(&c2f, "c2f", "", "chunk2file")
	flag.StringVar(&o2c, "o2c", "", "oid2locsfile")
	flag.StringVar(&outdir, "of", "", "output directory")
	flag.Parse()

	oidObjMap := getOidObjectMap(obf)
	fmt.Println("NUM OF OBJECTS: ", len(oidObjMap))
	chunkLocMap := getChunk2FileMap(c2f)
	fmt.Println("NUM OF CHUNKS : ", len(chunkLocMap))

	var count int
	var numobjs int
	for e := range getOid2Chunk(o2c) {
		count++
		if count%100000 == 0 {
			fmt.Println("Processed ", count)
		}
		if objname, ok := oidObjMap[e.oid]; ok {
			//if len(e.locs) != 1 || len(e.indices) != 1 {
			if len(e.locs) != 1 {
				fmt.Println("## SPECIAL HANDLE: ", objname)
			}
			if cloc, ok := chunkLocMap[e.locs[0].chunkId]; ok && e.locs[0].startOffset >= cloc.validRelativeOffset {
				fmt.Println("FOUND: ", objname, cloc)
				ifile, err := os.OpenFile(cloc.fname, os.O_RDONLY, 0444)
				checkErr(err)

				ifile.Seek(cloc.startOffset+e.locs[0].startOffset, os.SEEK_SET)

				var size int32
				if err := binary.Read(ifile, binary.BigEndian, &size); err != nil || int64(size) != e.locs[0].endOffset-e.locs[0].startOffset-48 {
					ifile.Close()
					fmt.Printf(" invalid block length: expected %d, but got %d\n", e.locs[0].endOffset-e.locs[0].startOffset-48, size)
					continue
				}

				chunkId := make([]byte, storage.ChunkIDLen)
				if n, err := ifile.ReadAt(chunkId, cloc.startOffset+e.locs[0].startOffset+4+int64(size)+storage.CRCLen); err != nil || n < storage.ChunkIDLen {
					ifile.Close()
					fmt.Println(" invalid chunk id")
					continue
				}

				if string(chunkId) != e.locs[0].chunkId {
					fmt.Printf(" invalid chunk id: expected %s, but got %s\n", e.locs[0].chunkId, string(chunkId))
					ifile.Close()
					continue
				}

				ifile.Seek(cloc.startOffset+e.locs[0].startOffset+4, os.SEEK_SET)

				targetdir := outdir + "." + strconv.Itoa(numobjs/10000)
				if _, err := os.Stat(targetdir); os.IsNotExist(err) {
					os.Mkdir(targetdir, os.ModePerm)
				}

				ofile, err := os.Create(targetdir + "/" + objname)
				checkErr(err)

				if _, err := io.CopyN(ofile, ifile, int64(size)); err != nil {
					fmt.Println("## Copy Error :", objname, err)
				}
				numobjs++
				fmt.Println("COPY DONE: ", targetdir+"/"+objname)

				ofile.Close()
				ifile.Close()
			}
		}
	}
}

func getOidObjectMap(fn string) map[string]string {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0444)
	checkErr(err)
	defer f.Close()

	r := bufio.NewReader(f)
	m := make(map[string]string)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			break
		}
		s = strings.TrimSpace(strings.TrimRight(s, "\n"))

		sum := sha256.Sum256([]byte("402c22d5fcfb78d2a7395ff9e353a65449b6bc5b5dcce62b73269ff49c0e341a." + s))
		m[fmt.Sprintf("%x", sum)] = s
	}
	return m
}

type ChunkLoc struct {
	fname               string
	startOffset         int64
	validRelativeOffset int64
	relativeEndOffset   int64
}

func getChunk2FileMap(fn string) map[string]*ChunkLoc {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0444)
	checkErr(err)
	defer f.Close()

	r := bufio.NewReader(f)
	m := make(map[string]*ChunkLoc)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			break
		}
		s = strings.TrimRight(s, "\n")
		parts := strings.Split(s, " ")
		// chunkid, filename, startOffset, validOffset, ChunkSize
		if len(parts) == 7 {
			so, _ := strconv.ParseInt(strings.TrimSpace(parts[4]), 10, 64)
			vro, _ := strconv.ParseInt(strings.TrimSpace(parts[5]), 10, 64)
			ro, _ := strconv.ParseInt(strings.TrimSpace(parts[6]), 10, 64)
			m[strings.TrimSpace(parts[2])] = &ChunkLoc{strings.TrimSpace(parts[3]), so, vro, ro}
		}
	}
	return m
}

type loc struct {
	chunkId     string
	startOffset int64
	endOffset   int64
}

type index struct {
	offset int64
	length int64
}

type OidLocs struct {
	oid     string
	locs    []*loc
	indices []*index
}

// 2017/07/14 17:20:09 27dfb875965af8f5eda9ab3247533d34dc0e263aa723e0ce02c3ba3d8f245fa2 - NLOC 1 loc_0 44b21c59-4249-437f-87f7-4a03f9c3bb28 16601959 16765728 163721 163721 NIDX 1 idx_0 0 163721
func parseEntry(s string) *OidLocs {
	parts := strings.Split(s, " - NLOC ")
	// 2017/07/14 17:20:09 27dfb875965af8f5eda9ab3247533d34dc0e263aa723e0ce02c3ba3d8f245fa2
	// 27dfb875965af8f5eda9ab3247533d34dc0e263aa723e0ce02c3ba3d8f245fa2
	oid := strings.Split(parts[0], " ")[2]
	// 1 loc_0 44b21c59-4249-437f-87f7-4a03f9c3bb28 16601959 16765728 163721 163721 NIDX 1 idx_0 0 163721
	locsindices := strings.Split(parts[1], " NIDX ")
	// 1 loc_0 44b21c59-4249-437f-87f7-4a03f9c3bb28 16601959 16765728 163721 163721
	locs := locsindices[0]
	// 1 idx_0 0 163721
	var indices string
	if len(locsindices) > 1 {
		indices = locsindices[1]
	}

	res := &OidLocs{oid: oid}

	nlocs := strings.SplitN(locs, " ", 2)[0]
	nl, _ := strconv.Atoi(nlocs)
	// loc_0 44b21c59-4249-437f-87f7-4a03f9c3bb28 16601959 16765728 163721 163721
	fields := strings.Split(strings.SplitN(locs, " ", 2)[1], " ")
	for i := 0; i < nl; i++ {
		id := fields[i*6+1]
		so, _ := strconv.ParseInt(fields[i*6+2], 10, 64)
		eo, _ := strconv.ParseInt(fields[i*6+3], 10, 64)
		res.locs = append(res.locs, &loc{id, so, eo})
	}

	if indices != "" {
		nindices := strings.SplitN(indices, " ", 2)[0]
		nidx, _ := strconv.Atoi(nindices)
		fields = strings.Split(strings.SplitN(indices, " ", 2)[1], " ")
		for i := 0; i < nidx; i++ {
			of, _ := strconv.ParseInt(fields[i*3+1], 10, 64)
			l, _ := strconv.ParseInt(fields[i*3+2], 10, 64)
			res.indices = append(res.indices, &index{of, l})
		}
	}

	return res
}

func getOid2Chunk(fn string) <-chan *OidLocs {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0444)
	checkErr(err)
	r := bufio.NewReader(f)

	c := make(chan *OidLocs)

	go func() {
		defer close(c)
		defer f.Close()

		for {
			s, err := r.ReadString('\n')
			if err != nil {
				break
			}
			s = strings.TrimRight(s, "\n")

			if fx, err := os.OpenFile(s, os.O_RDONLY, 0444); err == nil {
				rx := bufio.NewReader(fx)
				for {
					entry, err := rx.ReadString('\n')
					if err != nil {
						break
					}
					entry = strings.TrimRight(entry, "\n")
					if e := parseEntry(entry); e != nil {
						c <- e
					}
				}
				fx.Close()
			}
		}
	}()

	return c
}
