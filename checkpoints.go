package tensorflow_common

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetHighestNumberedCheckpoint will return the highest numbered checkpoint in a directory (non-recursively).  It
// accepts an optional "prefix", which is the filename prefix for the checkpoint files, defaulting to "model.ckpt-".
func GetHighestNumberedCheckpoint(directory string, prefix ...string) (int, error) {
	checkpointPrefix := "model.ckpt-"
	if len(prefix) > 0 {
		checkpointPrefix = prefix[0]
	}

	highestCheckpoint := 0
	if dirContents, err := ioutil.ReadDir(directory); err != nil {
		return 0, err
	} else {
		for _, dirContent := range dirContents {
			if strings.HasPrefix(dirContent.Name(), checkpointPrefix) && strings.HasSuffix(dirContent.Name(), ".index") {
				t, _ := strconv.Atoi(strings.Split(strings.Split(dirContent.Name(), "-")[1], ".")[0])

				if highestCheckpoint < t {
					highestCheckpoint = t
				}
			}
		}
	}

	return highestCheckpoint, nil
}