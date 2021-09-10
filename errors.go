package main

import (
	"github.com/pkg/errors"
)

var (
	ErrorFileNotFound       = errors.New("file not found")
	ErrorUnsupportedArchive = errors.New("archive type not supported")
	ErrorDownload           = errors.New("cannot download file")
	ErrorBadFile            = errors.New("cannot handle file inside archive")
)
