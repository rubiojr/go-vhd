package main

import (
	"./vhd"
	"fmt"
	"github.com/dustin/go-humanize"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: gh-vhd create <file> <size-in-bytes>\n")
		os.Exit(1)
	}

	file := os.Args[1]
	size := os.Args[2]

	isize, err := humanize.ParseBytes(size)

	if err != nil {
		panic(err)
	}

	vhd.CreateSparseVHD(uint64(isize), file)
	fmt.Printf("File %s (%s) created\n", file, humanize.IBytes(uint64(isize)))
}
