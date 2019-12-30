package tensorflow_common

import (
	"bufio"
	"os"
)

// LoadLabelsFileOnePerLine loads labels from a file, such that they are ordered and each is on a line by itself
func LoadLabelsFileOnePerLine(labelsFile string) ([]string, error) {
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
