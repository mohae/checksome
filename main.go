package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	input     string
	output    string
	uri       string
	checksum  string
	chunkSize int
	verbose   bool
)

func init() {
	flag.StringVar(&input, "input", "stdin", "input source")
	flag.StringVar(&input, "i", "stdin", "input source (short)")
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&uri, "url", "", "url of the input, takes precedence over -input")
	flag.StringVar(&uri, "u", "", "url of the input, takes precedence over -input (short)")
	flag.IntVar(&chunkSize, "chunksize", 8192, "size, in bytes, of each read from the input")
	flag.IntVar(&chunkSize, "s", 8192, "size, in bytes, of each read from the input (short)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&verbose, "v", false, "verbose output (short)")
	flag.StringVar(&checksum, "checksum", "sha256", "checksum algorithm")
	flag.StringVar(&checksum, "c", "sha256", "checksum algorithm (short)")

}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()
	typ, err := checksumFromString(checksum)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// set to defaults
	var (
		in, out *os.File
		n       int64
	)
	// set input
	in = os.Stdin
	if input != "stdin" {
		in, err = os.Open(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	defer in.Close()
	// set output
	out = os.Stdout
	if output != "stdout" {
		//
		out, err = os.OpenFile(output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		defer out.Close()
	}

	n, err = calcSum(typ, chunkSize, in, out)
	if err != nil {
		log.Printf("error calculating %s: %s", typ, err)
		return 1
	}
	out.WriteString("\n")
	if verbose {
		fmt.Printf("%s of %s: %d bytes read\n", typ, input, n)
	}
	return 0
}
