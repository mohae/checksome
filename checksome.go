// Copyright 2016 Joel Scoble
// Licensed under the MIT License;
// you may not use this file except in compliance with the License.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package checksome

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

const defaultChunk = 4096

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

var Upper bool

// ChecksumFromString returns the Checksum from the provided string.  For
// comparisons, hash names are normalized to lower case
func ChecksumFromString(name string) (Checksum, error) {
	switch strings.ToLower(name) {
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
		return Unknown, fmt.Errorf("Checksum: unsupported hash function type: %s", name)
	}
}

// GetHasher returns the hash.Hash for the supplied Checksum.  If the Checksum
// is Unknown, a nil will be returned.
func GetHasher(typ Checksum) hash.Hash {
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

// CalcSum takes a Checksum type, buffer size (chunk), reader, and writer;
// reading the data from reader and writing the resulting checksum to the
// writer.  The Checksum specifies the hash function to use for the
// calculation.  The number of bytes read and any error encountered is
// returned.
func CalcSum(c Checksum, chunk int, r io.Reader, w io.Writer) (n int64, err error) {
	if chunk < 1 {
		return 0, fmt.Errorf("invalid chunk size: %d", chunk)
	}
	h := GetHasher(c)
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
	if Upper {
		fmt.Fprintf(w, "%X", bs)
	} else {
		fmt.Fprintf(w, "%x", bs)
	}
	return n, nil
}

// Sum calculates the checksum of the data using the specified hash function.
func Sum(typ Checksum, data []byte) []byte {
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

// HashSHA256 calculates hashes using the SHA256 hashing function.
type HashSHA256 struct {
	r    io.Reader
	w    io.Writer
	h    hash.Hash
	buff *bufio.Reader
}

// TODO: should this be interface based (probably)

// NewSHA256 returns a *HashSHA256 that uses the provided io.Reader and
// io.Writer.  If bufferSize > 0; the read buffer will be sized to that value.
func NewSHA256(r io.Reader, w io.Writer, bufferSize int) *HashSHA256 {
	h := HashSHA256{
		r: r,
		w: w,
		h: sha256.New(),
	}
	if bufferSize > 0 {
		h.buff = bufio.NewReaderSize(r, chunk)
	} else {
		h.buff = bufio.NewReaderSize(r, defaultChunk)
	}
	return &h
}

// Sum calculates the hash of the bytes in the reader and writes the value to
// the writer.  The number of bytes read is returned.  If there is an error,
// the error is returned along with the number of bytes successfully read.
func (h *HashSHA256) Sum() (n int64, err error) {
	// refactor to use a generic implementation
	var x int64
	for {
		x, err = io.Copy(h.h, h, buff)
		n += x
		if err != nil {
			if err == io.EOF {
				break
			}
			return h.n, err
		}
		// if 0 bytes were read; at end
	}
}
