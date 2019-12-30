package tensorflow_common

import (
	"bufio"
	"os"
)

func LoadLabelsOnePerLine(labelsFile string) ([]string, error) {
	file, err := os.Open(labelsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return labels, nil
}
