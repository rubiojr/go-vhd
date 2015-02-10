package main

import (
	"./vhd"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"os"
	"strconv"
)

func createVHD(file, size string, options vhd.VHDOptions) {

	isize, err := humanize.ParseBytes(size)

	if err != nil {
		panic(err)
	}

	vhd.VHDCreateSparse(uint64(isize), file, options)
	fmt.Printf("File %s (%s) created\n", file, humanize.IBytes(uint64(isize)))
}

func vhdInfo(vhdFile string) {

	f, err := os.Open(vhdFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	vhd := vhd.FromFile(f)
	vhd.PrintInfo()
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

				opts := vhd.VHDOptions{}

				tstamp := c.String("timestamp")
				if tstamp != "" {
					itstamp, err := strconv.Atoi(tstamp)
					if err != nil {
						panic(err)
					}
					opts.Timestamp = int64(itstamp)
				}

				uuid := c.String("uuid")
				if uuid != "" {
					opts.UUID = uuid
				}
				createVHD(c.Args()[0], c.Args()[1], opts)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uuid",
					Value: "",
					Usage: "Set the UUID of the VHD header",
				},
				cli.StringFlag{
					Name:  "timestamp",
					Value: "",
					Usage: "Set the timestamp of the VHD header (UNIX time format)",
				},
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
