package tensorflow_common

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Decompress extracts a tar.gz or zip file and writes decompressed files to disk
func Decompress(pathToCompressedFile string) error {
	if strings.HasSuffix(pathToCompressedFile, ".tar.gz") {
		reader, err := os.Open(pathToCompressedFile)
		if err != nil {
			return err
		}
		defer reader.Close()

		r, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer r.Close()

		tarReader := tar.NewReader(r)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if err := processTarRecord(pathToCompressedFile, tarReader, header); err != nil {
				return err
			}
		}

		return nil
	}

	// ZIP
	r, err := zip.OpenReader(pathToCompressedFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		src, err := f.Open()
		if err != nil {
			return err
		}
		log.Println("Extracting", f.Name)
		dst, err := os.OpenFile(filepath.Join(path.Dir(pathToCompressedFile), f.Name), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		dst.Close()
	}
	return nil
}

func processTarRecord(
	pathToCompressedFile string,
	tarReader *tar.Reader,
	header *tar.Header,
) error {
	path := filepath.Join(path.Dir(pathToCompressedFile), header.Name)
	info := header.FileInfo()
	if info.IsDir() {
		if err := os.MkdirAll(path, info.Mode()); err != nil {
			return err
		}
		return nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, tarReader)
	if err != nil {
		return err
	}
	return nil
}
