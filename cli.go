package main

import (
	"./vhd"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"os"
)

func createVHD(file, size string) {

	isize, err := humanize.ParseBytes(size)

	if err != nil {
		panic(err)
	}

	vhd.CreateSparseVHD(uint64(isize), file)
	fmt.Printf("File %s (%s) created\n", file, humanize.IBytes(uint64(isize)))
}

func vhdInfo(vhdFile string) {

	f, err := os.Open(vhdFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	vhd.PrintVHDHeaders(f)
}

func main() {
	app := cli.NewApp()
	app.Name = "govhd"
	app.Usage = "Library and tool to manipulate VHD images"
	app.Author = "Sergio Rubio"
	app.Email = "rubiojr@frameos.org"

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "Create a VHD",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 2 {
					println("Missing command arguments.\n")
					fmt.Printf("Usage: %s create <file-path> <size MiB|GiB|...>\n",
						app.Name)
					os.Exit(1)
				}
				createVHD(c.Args()[0], c.Args()[1])
			},
		},
		{
			Name:  "info",
			Usage: "Print VHD info",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					println("Missing command arguments.\n")
					fmt.Printf("Usage: %s info <file-path> <size MiB|GiB|...>\n",
						app.Name)
					os.Exit(1)
				}
				vhdInfo(c.Args()[0])
			},
		},
	}

	app.Run(os.Args)
}
