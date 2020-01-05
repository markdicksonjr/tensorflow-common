package tensorflow_common

import (
	"bufio"
	"os"
	"strings"
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
		text := strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			labels = append(labels, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return labels, nil
}
