package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Archive struct {
	URL string
}

// ExtractFile Returns the bytes of the file with `name` from the archive.
func (a *Archive) ExtractFile(name string) ([]byte, error) {
	r, err := wget(a.URL)
	if err != nil {
		return []byte{}, ErrorDownload
	}

	switch {
	case strings.HasSuffix(a.URL, ".tar.gz"):
		return extractTar(r, name)
	case strings.HasSuffix(a.URL, ".zip"):
		return extractZip(r, name)
	}

	return []byte{}, ErrorUnsupportedArchive
}

func extractZip(data []byte, name string) ([]byte, error) {
	zipr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return []byte{}, err
	}

	for _, f := range zipr.File {
		if f.Name != name {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return []byte{}, errors.Wrap(err, "cannot open file")
		}

		file, err := ioutil.ReadAll(rc)
		return file, err
	}

	return []byte{}, ErrorFileNotFound
}

func extractTar(data []byte, name string) ([]byte, error) {
	uncompressed, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}

	tarReader := tar.NewReader(uncompressed)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return []byte{}, err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			headerPath := stripDir(header.Name)

			if headerPath != name {
				continue
			}

			file, err := ioutil.ReadAll(tarReader)

			log.Println("found", header.Name)

			return file, err
		default:
			errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %+v in %s",
				header.Typeflag,
				header.Name))
		}
	}

	return []byte{}, ErrorFileNotFound
}

// stripDir Remove the first directory of a file path. To be used when
// searching for files inside a zip/tar file as the first part is the zip name.
func stripDir(path string) string {
	fields := strings.Split(path, "/")

	return filepath.Join(fields[1:len(fields)]...)
}

func wget(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot download file")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	return ioutil.ReadAll(response.Body)
}
