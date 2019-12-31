package tensorflow_common

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"log"
	"sort"
)

type ObjectDetectionResult struct {
	Probability float32
	Box         []float32
	Label       string
	Index       int
	Class       float32
}

type ObjectDetectionOperationMappings struct {
	ImageTensor      string
	DetectionBoxes   string
	DetectionScores  string
	DetectionClasses string
	NumDetections    string
}

// RawSlicesToObjectDetectionResult will convert the slices that come from inference to more convenient results
func RawSlicesToObjectDetectionResult(
	probabilities []float32,
	classes []float32,
	boxes [][]float32,
	labels []string,
) []ObjectDetectionResult {
	var results []ObjectDetectionResult

	for i, p := range probabilities {
		labelIndex := int(classes[i])
		if labelIndex > len(labels) {
			labelIndex = 0
		}

		results = append(results, ObjectDetectionResult{
			Probability: p,
			Box:         boxes[i],
			Label:       labels[labelIndex],
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

// PerformObjectDetectionInference runs the image through the model
func PerformObjectDetectionInference(
	input *tf.Tensor,
	graph *tf.Graph,
	session *tf.Session,
	mappings ...ObjectDetectionOperationMappings,
) (probabilities, classes []float32, boxes [][]float32, err error) {
	var resolvedMappings ObjectDetectionOperationMappings

	if len(mappings) > 0 {
		resolvedMappings = mappings[0]
	} else {
		resolvedMappings = ObjectDetectionOperationMappings{
			ImageTensor:      "image_tensor",
			DetectionBoxes:   "detection_boxes",
			DetectionScores:  "detection_scores",
			DetectionClasses: "detection_classes",
			NumDetections:    "num_detections",
		}
	}

	// get all the input and output operations
	inputOp := graph.Operation(resolvedMappings.ImageTensor)
	ops := graph.Operations()
	for _, v := range ops {
		log.Println(v.Name())
	}

	// get the output ops
	o1 := graph.Operation(resolvedMappings.DetectionBoxes)
	o2 := graph.Operation(resolvedMappings.DetectionScores)
	o3 := graph.Operation(resolvedMappings.DetectionClasses)
	o4 := graph.Operation(resolvedMappings.NumDetections)

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			inputOp.Output(0): input,
		},
		[]tf.Output{
			o1.Output(0),
			o2.Output(0),
			o3.Output(0),
			o4.Output(0),
		},
		nil)
	if err != nil {
		return nil, nil, nil, err
	}

	allProbabilities := output[1].Value().([][]float32)
	probabilities = allProbabilities[0]
	classes = output[2].Value().([][]float32)[0]
	boxes = output[0].Value().([][][]float32)[0]

	return
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
