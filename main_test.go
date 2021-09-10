package main

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFile(t *testing.T) {
	var cases = []struct {
		name    string
		archive *Archive
		sha256  string
		file    string
		err     error
	}{
		{
			name: "binary from URL zip",
			archive: &Archive{
				URL: "https://github.com/terraform-linters/tflint/releases/download/v0.31.0/tflint_darwin_amd64.zip",
			},
			file:   "tflint",
			sha256: "0d226ac5664393b8ce088c8fb4275aa2dcefdf767a12e15030339ec53184e5d0",
		},
		{
			name: "binary from URL tar.gz",
			archive: &Archive{
				URL: "https://github.com/cli/cli/releases/download/v2.0.0/gh_2.0.0_linux_386.tar.gz",
			},
			sha256: "8c8a28c1fdd17b9e2b21892c760d344e4d53ca225dbe9b9fcf9404cecfdd19f0",
			file:   "bin/gh",
		},
		{
			name: "file not found",
			archive: &Archive{
				URL: "https://github.com/cli/cli/releases/download/v2.0.0/gh_2.0.0_linux_386.tar.gz",
			},
			file:   "bin/gh1",
			sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			err:    ErrorFileNotFound,
		},
		{
			name: "unsupported archive type",
			archive: &Archive{
				URL: "https://github.com/mhristof.png",
			},
			sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			err:    ErrorUnsupportedArchive,
		},
		{
			name: "wrong url",
			archive: &Archive{
				URL: "https://github.com/mhristof1111111.png",
			},
			sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			err:    ErrorDownload,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fmt.Println(fmt.Sprintf("test.archive: %+v", test.archive))
			data, err := test.archive.ExtractFile(test.file)
			assert.ErrorIs(t, test.err, err, test.name)
			h := sha256.New()
			h.Write(data)
			assert.Equal(t, test.sha256, fmt.Sprintf("%x", h.Sum(nil)), test.name)
		})
	}
}

func TestStripDir(t *testing.T) {
	var cases = []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "multilevel path",
			in:   "foo/bar/baz",
			out:  "bar/baz",
		},
		{
			name: "single file",
			in:   "foo",
			out:  "",
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.out, stripDir(test.in), test.name)
	}
}
