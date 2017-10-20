package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Usage: gzdiff base.gz newer.gz")
		os.Exit(1)
	}
	bpath := os.Args[1]
	npath := os.Args[2]
	ext := filepath.Ext(npath)
	dpath := strings.TrimSuffix(npath, ext) + ".diff" + ext

	lh := map[uint64]bool{}
	gzipEachLine(bpath, func(line string, hash uint64) {
		lh[hash] = true
	})

	f, _ := os.Create(dpath)
	w := gzip.NewWriter(f)
	defer w.Close()

	gzipEachLine(npath, func(line string, hash uint64) {
		if !lh[hash] {
			w.Write([]byte(line + "\n"))
		}
	})
}

func gzipEachLine(path string, proc func(line string, hash uint64)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		h := xxhash.Sum64String(l)
		proc(l, h)
	}
	return nil
}
