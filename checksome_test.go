package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestProcessChecksumTypes(t *testing.T) {
	tests := []struct {
		vals     []string
		expected []Checksum
		err      string
	}{
		{[]string{"SHA1"}, []Checksum{SHA1}, ""},
		{
			[]string{"SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA512_224", "SHA512_256"},
			[]Checksum{SHA1, SHA224, SHA256, SHA384, SHA512, SHA512_224, SHA512_256},
			"",
		},
		{
			[]string{"SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA512_224", "SHA512_256"},
			[]Checksum{SHA1, SHA224, SHA256, SHA384, SHA512, SHA512_224, SHA512_256},
			"",
		},
		{
			[]string{"SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA512_224", "SHA512_256"},
			[]Checksum{SHA1, SHA224, SHA256, SHA384, SHA512, SHA512_224, SHA512_256},
			"",
		},
		{[]string{"wut"}, []Checksum{}, "unknown checksum type: wut"},
	}
	for i, test := range tests {
		res, err := processChecksumTypes(test.vals)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%d: got %q want %q", i, err, test.err)
			}
			continue
		}
		if test.err != "" {
			t.Errorf("%d: got no error, expected %q", i, test.err)
			continue
		}
		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("%d: got %v want %v", i, res, test.expected)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		typ      Checksum
		expected []byte
	}{
		{SHA1, []byte{0x17, 0x9d, 0xf6, 0x01, 0xa4, 0xe8, 0xcf, 0xd0, 0xa9, 0x7b, 0x29, 0x48, 0x44, 0x10, 0xa5, 0xf4, 0x2f, 0xcb, 0xff, 0xf1}},
		{SHA224, []byte{0x8f, 0x6f, 0x10, 0x10, 0x84, 0x58, 0x8d, 0x8a, 0x6e, 0xd3, 0xe6, 0x25, 0x99, 0xac, 0x6b, 0xe9, 0xfd, 0x1f, 0x6b, 0x7c, 0xa4, 0x43, 0x18, 0x9a, 0x29, 0x77, 0x87, 0x81}},
		{SHA256, []byte{0x02, 0xb5, 0xdc, 0xd5, 0xf0, 0xef, 0x1a, 0x39, 0xcf, 0xfe, 0xc5, 0xf8, 0xb6, 0x25, 0xec, 0x20, 0xbf, 0xfc, 0xf9, 0x18, 0xe4, 0xef, 0xd3, 0xf5, 0x4b, 0xab, 0xec, 0x4e, 0xae, 0x03, 0xb3, 0x47}},
		{SHA384, []byte{0x4a, 0x90, 0xfc, 0x2b, 0xe8, 0x5e, 0xf1, 0xbf, 0x0d, 0x75, 0xb9, 0x19, 0x20, 0xd8, 0x56, 0x9a, 0x79, 0xcf, 0x62, 0xbc, 0x21, 0xef, 0x86, 0x6d, 0xf8, 0x9e, 0x31, 0xdf, 0x6f, 0x57, 0xda, 0xb5, 0x03, 0xbd, 0x2b, 0x70, 0x54, 0x85, 0xf1, 0x09, 0xe7, 0xde, 0xfb, 0x05, 0x4e, 0x87, 0xcd, 0x3d}},
		{SHA512, []byte{0x8c, 0xb2, 0xdd, 0x06, 0x29, 0x85, 0x9c, 0x90, 0x40, 0x65, 0x70, 0x0a, 0x76, 0x1c, 0xf2, 0x61, 0x6a, 0x26, 0xd8, 0xd4, 0x77, 0x56, 0x0e, 0xed, 0x47, 0x3d, 0x34, 0x89, 0xc5, 0x37, 0x20, 0x46, 0x15, 0x11, 0xb6, 0x6f, 0x19, 0x1d, 0x06, 0x53, 0x7e, 0x38, 0x87, 0xf4, 0xe6, 0xf4, 0x0c, 0xb9, 0xb0, 0x23, 0x12, 0x85, 0xb9, 0xb5, 0xdd, 0x69, 0x95, 0xf0, 0x82, 0x5b, 0x18, 0xd5, 0xf1, 0xaf}},
		{SHA512_224, []byte{0x53, 0x56, 0x4f, 0xa8, 0x2b, 0x93, 0xc6, 0xe9, 0xd0, 0xd2, 0x94, 0x70, 0x03, 0x8f, 0x51, 0x52, 0x24, 0x03, 0x6e, 0xf5, 0xe8, 0x7b, 0x84, 0xef, 0xb9, 0xdf, 0x5d, 0x4c}},
		{SHA512_256, []byte{0xcd, 0x45, 0x92, 0x9a, 0x5e, 0xe9, 0xd8, 0xe7, 0x49, 0x40, 0x9b, 0x8a, 0x08, 0x6a, 0x60, 0xd6, 0x7a, 0x30, 0x18, 0xe1, 0xe2, 0x7b, 0x9e, 0xfd, 0x8e, 0x76, 0xdb, 0x49, 0xd5, 0xd2, 0x68, 0x11}},
		{Unknown, nil},
	}
	data := []byte("Hello, World.")

	for _, test := range tests {
		h := sum(test.typ, data)
		if bytes.Compare(h, test.expected) != 0 {
			t.Errorf("%s: got %v want %v", test.typ, h, test.expected)
		}
	}
}

func TestCalcHash(t *testing.T) {
	tests := []struct {
		name     string
		typ      Checksum
		chunk    int
		expected string
		err      string
	}{
		{"SHA1", SHA1, 1024, "179df601a4e8cfd0a97b29484410a5f42fcbfff1", ""},
		{"SHA224", SHA224, 1024, "8f6f101084588d8a6ed3e62599ac6be9fd1f6b7ca443189a29778781", ""},
		{"SHA256", SHA256, 1024, "02b5dcd5f0ef1a39cffec5f8b625ec20bffcf918e4efd3f54babec4eae03b347", ""},
		{"SHA384", SHA384, 1024, "4a90fc2be85ef1bf0d75b91920d8569a79cf62bc21ef866df89e31df6f57dab503bd2b705485f109e7defb054e87cd3d", ""},
		{"SHA512", SHA512, 1024, "8cb2dd0629859c904065700a761cf2616a26d8d477560eed473d3489c53720461511b66f191d06537e3887f4e6f40cb9b0231285b9b5dd6995f0825b18d5f1af", ""},
		{"SHA512_224", SHA512_224, 1024, "53564fa82b93c6e9d0d29470038f515224036ef5e87b84efb9df5d4c", ""},
		{"SHA512_256", SHA512_256, 1024, "cd45929a5ee9d8e749409b8a086a60d67a3018e1e27b9efd8e76db49d5d26811", ""},
		{"no hasher error", Unknown, 1024, "", "unknown checksum type: Unknown"},
		{"no buffer error", SHA1, 0, "", "invalid chunk size: 0"},
	}

	data := []byte("Hello, World.")

	for _, test := range tests {
		b := make([]byte, 0, 64)
		r := bytes.NewBuffer(data)
		w := bytes.NewBuffer(b)
		n, err := calcSum(test.typ, test.chunk, r, w)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%s: got %q, want %q", test.name, err, test.err)
			}
		} else {
			if test.err != "" {
				t.Errorf("%s: got no error, want %q", test.name, test.err)
				continue
			}
			if n != 13 {
				t.Errorf("%s: got %d bytes read, want 13", test.name, n)
			}
			if test.expected != w.String() {
				t.Errorf("%s: got %s want %s", test.name, w, test.expected)
			}
		}
	}
}
