package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	cdfPrefix  string
	blobPrefix string
	outdir     string
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&cdfPrefix, "cdf", "result.cdf.", "cdf prefix")
	flag.StringVar(&blobPrefix, "blob", "result.blob.", "blob prefix")
	flag.StringVar(&outdir, "out", "", "output directory")
	flag.Parse()

	blobs := getBlobs(blobPrefix)
	fmt.Println("NUM BLOBS ", len(blobs))

	files, err := filepath.Glob(cdfPrefix + "*/*")
	checkErr(err)
	var count, numobjs, idx int
	var cdfname, msg string
	for _, f := range files {
		count++
		if count%100000 == 0 {
			fmt.Println("Processed ", count)
		}
		cdfname = filepath.Base(f)
		ifile, err := os.OpenFile(f, os.O_RDONLY, 0444)
		checkErr(err)
		body, err := ioutil.ReadAll(ifile)
		checkErr(err)
		ifile.Close()
		msg = string(body)
		if idx = strings.Index(msg, "eclipblob md5="); idx == -1 {
			continue
		}
		msg = msg[idx:]
		parts := strings.SplitN(msg, "\"", 5)
		if len(parts) != 5 {
			continue
		}
		blobname, size := parts[1], parts[3]
		if _, ok := blobs[blobname]; !ok {
			continue
		}
		if stats, err := os.Stat(blobs[blobname]); err != nil || strconv.FormatInt(stats.Size(), 10) != size {
			continue
		}

		targetdir := outdir + "." + strconv.Itoa(numobjs/10000)
		if _, err := os.Stat(targetdir); os.IsNotExist(err) {
			os.Mkdir(targetdir, os.ModePerm)
		}
		fmt.Println(f, blobs[blobname], targetdir+"/"+cdfname, targetdir+"/"+cdfname+"_blob")
		os.Rename(f, targetdir+"/"+cdfname)
		os.Rename(blobs[blobname], targetdir+"/"+cdfname+"_blob")
		numobjs++
	}
}

func getBlobs(prefix string) map[string]string {
	files, err := filepath.Glob(prefix + "*/*")
	checkErr(err)
	m := make(map[string]string)
	for _, f := range files {
		m[filepath.Base(f)] = f
	}
	return m
}
