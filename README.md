# Tensorflow Common

Utilities for the Go Tensorflow API

## Setup

https://www.tensorflow.org/install/lang_go

## Sample

```go
graph, err := tensorflow_playground.LoadGraphFromFile(pbFilePath)
if err != nil {
    log.Fatal(err)
}

// create a session for inference over graph.
session, err := tf.NewSession(graph, nil)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

// run inference on *imageFile.
// for multiple images, session.Run() can be called in a loop (and
// concurrently). Alternatively, images can be batched since the model
// accepts batches of image data as input.
tensor, image, err := tensorflow_common.LoadTensorFromImageFile(imageFilePath)
if err != nil {
    log.Fatal(err)
}

probabilities, classes, boxes, err := tensorflow_common.PerformObjectDetectionInference(tensor, graph, session)
if err != nil {
    log.Fatalf("error making prediction: %v", err)
}

labels, err := tensorflow_common.LoadLabelsFileOnePerLine(labelsFilePath)
if err != nil {
    log.Fatal(err)
}

results := tensorflow_common.RawSlicesToObjectDetectionResult(probabilities, classes, boxes, labels)
finalResults := tensorflow_common.FilterByMinConfidence(results, 0.8)
```

## Utilities

Among other things, this repo can convert from Pascal VOC to CSV.  This is handy for using LabelImg and working toward a
TFRecord for object detection.