package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"strings"
)

//go:generate stringer -type=Checksum
type Checksum int

// Possible Hashes
const (
	Unknown Checksum = iota
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	SHA512_224
	SHA512_256
)

// checksumFromString returns the Checksum from the provided string.
func checksumFromString(v string) (Checksum, error) {
	switch strings.ToLower(v) {
	case "sha1":
		return SHA1, nil
	case "sha224":
		return SHA224, nil
	case "sha256":
		return SHA256, nil
	case "sha384":
		return SHA384, nil
	case "sha512":
		return SHA512, nil
	case "sha512_224":
		return SHA512_224, nil
	case "sha512_256":
		return SHA512_256, nil
	default:
		return Unknown, fmt.Errorf("unknown checksum type: %s", v)
	}
}

func getHasher(typ Checksum) hash.Hash {
	switch typ {
	case SHA1:
		return sha1.New()
	case SHA224:
		return sha256.New224()
	case SHA256:
		return sha256.New()
	case SHA384:
		return sha512.New384()
	case SHA512:
		return sha512.New()
	case SHA512_224:
		return sha512.New512_224()
	case SHA512_256:
		return sha512.New512_256()
	}
	return nil
}

// calcSum takes a Checksum type, buffer size (chunk), reader, and writer;
// reading the data from reader and writing the resulting checksum to the
// writer.  The checksum type is used to determine which algorithm will be
// used.  The number of bytes read and any error encountered is returned.
func calcSum(c Checksum, chunk int, r io.Reader, w io.Writer) (n int64, err error) {
	if chunk < 1 {
		return 0, fmt.Errorf("invalid chunk size: %d", chunk)
	}
	h := getHasher(c)
	if h == nil {
		return 0, fmt.Errorf("unknown checksum type: %s", c)
	}
	buf := bufio.NewReaderSize(r, chunk)
	var x int64
	for {
		x, err = io.Copy(h, buf)
		n += x
		if err != nil {
			if err == io.EOF {
				break
			}
			return x, err
		}
		// if 0 bytes were read; at end
		if x == 0 {
			break
		}
	}
	bs := h.Sum(nil)
	if upper {
		fmt.Fprintf(w, "%X", bs)
	} else {
		fmt.Fprintf(w, "%x", bs)
	}
	return n, nil
}

func sum(typ Checksum, data []byte) []byte {
	switch typ {
	case SHA1:
		h := sha1.Sum(data)
		return h[:]
	case SHA224:
		h := sha256.Sum224(data)
		return h[:]
	case SHA256:
		h := sha256.Sum256(data)
		return h[:]
	case SHA384:
		h := sha512.Sum384(data)
		return h[:]
	case SHA512:
		h := sha512.Sum512(data)
		return h[:]
	case SHA512_224:
		h := sha512.Sum512_224(data)
		return h[:]
	case SHA512_256:
		h := sha512.Sum512_256(data)
		return h[:]
	}
	return nil
}
