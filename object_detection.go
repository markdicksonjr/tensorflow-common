package tensorflow_common

import "sort"

type ObjectDetectionResult struct {
	Probability float32
	Box         []float32
	Label       string
	Index       int
	Class       float32
}

func RawSlicesToObjectDetectionResult(
	probabilities []float32,
	classes []float32,
	boxes [][]float32,
	labels []string,
) []ObjectDetectionResult {
	var results []ObjectDetectionResult

	for i, p := range probabilities {
		results = append(results, ObjectDetectionResult{
			Probability: p,
			Box:         boxes[i],
			Label:       labels[int(classes[i])],
			Class:       classes[i],
			Index:       i,
		})
	}

	// sort the results so the highest probability is up top
	sort.Slice(results, func(i, j int) bool {
		return results[i].Probability > results[j].Probability
	})

	return results
}

func FilterByMinConfidence(
	results []ObjectDetectionResult,
	minConfidence float64,
) []ObjectDetectionResult {
	var finalResults []ObjectDetectionResult

	for _, p := range results {
		if float64(p.Probability) < minConfidence {
			continue
		}

		finalResults = append(finalResults, p)
	}

	return finalResults
}