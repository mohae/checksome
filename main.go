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
	chunkSize int
	checksums stringArray
)

type stringArray []string

func (f *stringArray) String() string {
	var s string
	for i, v := range *f {
		if i == 0 {
			s = v
			continue
		}
		s += ", " + v

	}
	return s
}

func (f *stringArray) Set(v string) error {
	*f = append(*f, v)
	return nil
}

func init() {
	flag.StringVar(&input, "input", "stdin", "input source")
	flag.StringVar(&input, "i", "stdin", "input source (short)")
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&uri, "url", "", "url of the input, takes precedence over -input")
	flag.StringVar(&uri, "u", "", "url of the input, takes precedence over -input (short)")
	flag.IntVar(&chunkSize, "chunksize", 8192, "size, in bytes, of each read from the input")
	flag.IntVar(&chunkSize, "s", 8192, "size, in bytes, of each read from the input (short)")
	flag.Var(&checksums, "checksum", "checksum algorithm to use: default is sha256")
	flag.Var(&checksums, "c", "checksum algorithm to use: default is sha256 (short)")

}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()
	if len(checksums) == 0 {
		checksums = append(checksums, "sha256")
	}
	typs, err := processChecksumTypes(checksums)
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

	// initially, only support calculation of 1 hash at a time.
	// mutli-hash support probably involves multi-writer
	n, err = calcSum(typs[0], chunkSize, in, out)
	if err != nil {
		log.Printf("error calculating %s: %s", typs[0], err)
		return 1
	}
	out.WriteString("\n")
	log.Printf("%s of %s calculated: %d bytes read", typs[0], input, n)
	return 0
}
