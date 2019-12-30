package tensorflow_common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func FilesExist(files ...string) bool {
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return false
		}
	}
	return true
}

func DirectoryExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func Download(URL, filename string) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func DownloadToBuffer(url string, cacheDir string) (buffer *bytes.Buffer, err error) {
	baseFilename := ""
	baseFilename, err = GetFilenameFromUrl(url)
	if err != nil {
		return
	}

	// if the file is in our cache, use the local cached copy
	if cacheDir != "" && FileExists(cacheDir+"/"+baseFilename) {
		if fileContents, readErr := ioutil.ReadFile(cacheDir + "/" + baseFilename); readErr != nil {
			err = readErr
			return
		} else {
			return bytes.NewBuffer(fileContents), nil
		}
	}

	// ignore certificate issues
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// get the data
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	b := bytes.Buffer{}
	_, err = io.Copy(&b, resp.Body)
	if err != nil {
		return &b, err
	}

	if cacheDir != "" {
		if err := ioutil.WriteFile(cacheDir+"/"+baseFilename, b.Bytes(), os.ModePerm); err != nil {
			return &b, err
		}
	}

	return &b, nil
}

func CopyFile(src, dst string) (int64, error) {
	src_file, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}

	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}

func GetFilenameFromUrl(baseModelUrl string) (string, error) {
	baseModelUrlParsed, err := url.Parse(baseModelUrl)
	if err != nil {
		return "", errors.Wrap(err, "could not parse URL")
	}
	baseModelUrlParts := strings.Split(baseModelUrlParsed.Path, "/")
	return baseModelUrlParts[len(baseModelUrlParts)-1], nil
}

// GraphFromFile reads the contents of a file and attempts to treat it as serialized
// graph data.  This is an easy way to get Tensorflow pb files into graphs
func GraphFromFile(pbFile string) (*tf.Graph, error) {
	modelFileContents, err := ioutil.ReadFile(pbFile)
	if err != nil {
		return nil, err
	}

	graph := tf.NewGraph()
	if err := graph.Import(modelFileContents, ""); err != nil {
		return graph, err
	}

	return graph, nil
}
