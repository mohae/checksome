// Copyright 2016 Joel Scoble
// Licensed under the MIT License;
// you may not use this file except in compliance with the License.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mohae/checksome"
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
	flag.StringVar(&checksum, "hash", "sha256", "hash algorithm")
	flag.StringVar(&checksum, "-s", "sha256", "hash algorithm (short)")
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
	// set whether or not the output should be uppercase
	checksome.Upper = upper
	typ, err := checksome.ChecksumFromString(checksum)
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
	n, err = checksome.CalcSum(typ, chunkSize, in, out)
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
