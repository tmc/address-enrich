package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tmc/addressenrich/usps"
)

func main() {
	// var flagInput = flag.String("input", "826 treat ave, san francisco", "input address")
	var flagInputFile = flag.String("input-file", "-", "input file")
	var flagVerbose = flag.Bool("v", false, "verbose mode")
	var flagSkipRows = flag.Int("skip-rows", 1, "rows to skip")
	var flagStartIndex = flag.Int("start-column", 0, "column start index")
	flag.Parse()

	os.Exit(run(RunOptions{
		InputFile:   *flagInputFile,
		Verbose:     *flagVerbose,
		SkipRows:    *flagSkipRows,
		StartColumn: *flagStartIndex,
	}))
}

type RunOptions struct {
	InputFile   string
	Verbose     bool
	SkipRows    int
	StartColumn int
}

func run(opts RunOptions) int {
	var file io.Reader
	var err error
	if opts.InputFile == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(opts.InputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	u := &usps.Client{
		Username:   os.Getenv("USPS_USERNAME"),
		Password:   os.Getenv("USPS_PASSWORD"),
		Production: true,
	}

	idx := func(s []string, i int) string {
		if i >= len(s) {
			return ""
		}
		return s[i]
	}

	scanner := bufio.NewScanner(file)
	rowIndex := 0
	for scanner.Scan() {
		rowIndex++
		if rowIndex <= opts.SkipRows {
			continue
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		// TODO: this is pretty hacky, it expects these fields to all be present
		a := usps.Address{
			Address1: idx(parts, 0+opts.StartColumn),
			City:     idx(parts, 1+opts.StartColumn),
			State:    idx(parts, 2+opts.StartColumn),
			Zip5:     idx(parts, 3+opts.StartColumn),
		}
		resp, err := u.ZipByAddress(a)
		if err != nil {
			fmt.Fprintln(os.Stderr, err, a.Address1)
		}
		address := resp.Address()
		if opts.Verbose {
			fmt.Printf("%+v\n", address)
		}
		fmt.Println(strings.Join([]string{
			line,
			address.AddressLine1,
		}, ","))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return -1
	}
	return 0
}
