package archive

import "github.com/pkg/errors"

var (
	// ErrorFileNotFound Error returned when the file requested is not found
	// in the archive.
	ErrorFileNotFound = errors.New("file not found")
	// ErrorUnsupportedArchive Error when the archive type is not supported.
	ErrorUnsupportedArchive = errors.New("archive type not supported")
	// ErrorDownload Error when there was an error while downloading the file.
	ErrorDownload = errors.New("cannot download file")
)
