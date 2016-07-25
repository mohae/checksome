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
	checksum  string
	chunkSize int
	help      bool
	upper     bool
	verbose   bool
)

func init() {
	flag.StringVar(&input, "input", "stdin", "input source")
	flag.StringVar(&input, "i", "stdin", "input source (short)")
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.IntVar(&chunkSize, "readchunk", 8192, "size, in bytes, of each read from the input")
	flag.IntVar(&chunkSize, "r", 8192, "size, in bytes, of each read from the input (short)")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&help, "h", false, "help (short)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&verbose, "v", false, "verbose output (short)")
	flag.BoolVar(&upper, "upper", false, "uppercase the output")
	flag.BoolVar(&upper, "u", false, "uppercase the output (short)")
	flag.StringVar(&checksum, "checksum", "sha256", "checksum algorithm")
	flag.StringVar(&checksum, "c", "sha256", "checksum algorithm (short)")
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()
	// the only arg we care about is help.  This is in case the user uses
	// just help instead of -help or -h
	for _, arg := range args {
		if arg == "help" {
			help = true
			break
		}
	}
	if help {
		Usage()
		return 1
	}
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
	if verbose {
		out.WriteString(fmt.Sprintf("%s:\t", typ.String()))
	}
	n, err = calcSum(typ, chunkSize, in, out)
	if err != nil {
		log.Printf("error calculating %s: %s", typ, err)
		return 1
	}
	out.WriteString("\n")
	if verbose {
		fmt.Printf("%s (%s): %d bytes read\n", input, typ, n)
	}
	return 0
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	flag.PrintDefaults()
}
