package main

import (
  "os"
	"fmt"
	"./vhd"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: vhd-reader <vhd-file>")
		os.Exit(0)
	}

	vhdFile := os.Args[1]
	f, err := os.Open(vhdFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("\nReading VHD header...")
	vhd.PrintVHDHeaders(f)

	// header should be equal to the footer, added for redundancy
	///fmt.Println("\nReading VHD footer...")
	///fstat, err := f.Stat()
	///check(err)
	///vhdFooter := make([]byte, 512)
	///f.Seek(fstat.Size()-512, 0)
	///_, err = f.Read(vhdFooter)
	///check(err)
	///readVHDHeader(vhdFooter)
}
