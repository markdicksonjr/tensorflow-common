package voc

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"strconv"
	"strings"
)

type Annotation struct {
	XMLName   xml.Name           `xml:"annotation"`
	Objects   []AnnotationObject `xml:"object"`
	Folder    string             `xml:"folder"`
	Filename  string             `xml:"filename"`
	Path      string             `xml:"path"`
	Source    AnnotationSource   `xml:"source"`
	Width     int                `xml:"size>width"`
	Height    int                `xml:"size>height"`
	Depth     int                `xml:"size>depth"`
	Segmented int                `xml:"segmented"`
}

type AnnotationObject struct {
	Name        string                `xml:"name"`
	Pose        string                `xml:"pose"`
	Truncated   int                   `xml:"truncated"`
	Difficult   int                   `xml:"difficult"`
	BoundingBox AnnotationBoundingBox `xml:"bndbox"`
}

type AnnotationBoundingBox struct {
	XMin int `xml:"xmin"`
	YMin int `xml:"ymin"`
	XMax int `xml:"xmax"`
	YMax int `xml:"ymax"`
}

type AnnotationSource struct {
	Database string `xml:"database"`
}

func ReadAnnotationsFromDir(path string) ([]Annotation, error) {
	var annotations []Annotation
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".xml") {
			annotation, err := ReadAnnotationFile(path + "/" + f.Name())
			if err != nil {
				return annotations, err
			}
			annotations = append(annotations, annotation)
		}
	}

	return annotations, nil
}

func ReadAnnotationFile(filename string) (Annotation, error) {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Annotation{}, err
	}

	return ReadAnnotation(fileContents)
}

func ReadAnnotation(contents []byte) (Annotation, error) {
	result := Annotation{}
	err := xml.Unmarshal(contents, &result)
	return result, err
}

func AnnotationsToCsv(annotations []Annotation) ([]byte, error) {
	rows := [][]string{
		{ "filename","width","height","class","xmin","ymin","xmax","ymax"},
	}

	for _, a := range annotations {
		for _, o := range a.Objects {
			rows = append(rows, []string{
				a.Filename, strconv.Itoa(a.Width),
				strconv.Itoa(a.Height),
				o.Name,
				strconv.Itoa(o.BoundingBox.XMin),
				strconv.Itoa(o.BoundingBox.YMin),
				strconv.Itoa(o.BoundingBox.XMax),
				strconv.Itoa(o.BoundingBox.YMax),
			})
		}
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	csvWriter := csv.NewWriter(bytesBuffer)
	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return nil, err
		}
	}

	csvWriter.Flush()
	return bytesBuffer.Bytes(), nil
}
