package tensorflow_common

import (
	"bytes"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

// convert the image in filename to a Tensor suitable as input to the Inception model.
func TensorFromImage(filename string) (*tf.Tensor, image.Image, error) {
	bytesVal, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytesVal))
	if err != nil {
		return nil, nil, err
	}

	// construct a graph to normalize the image
	graph, input, output, err := buildDecodeBitmapGraph()
	if err != nil {
		return nil, nil, err
	}

	// execute that graph to normalize this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	// run the session
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, nil, err
	}

	r := bytes.NewReader(bytesVal)

	if strings.HasSuffix(filename, "bmp") {
		i, _, err := image.Decode(r)
		if err != nil {
			return nil, nil, err
		}

		return normalized[0], i, nil
	} else if strings.HasSuffix(filename, "png") {

		i, err := png.Decode(r)
		if err != nil {
			return nil, nil, err
		}

		return normalized[0], i, nil
	} else if strings.HasSuffix(filename, "jpg") || strings.HasSuffix(filename, "jpeg") {
		i, err := jpeg.Decode(r)
		if err != nil {
			return nil, nil, err
		}

		return normalized[0], i, nil
	}

	return normalized[0], nil, nil
}

// build a graph to decode bitmap input into the proper tensor shape
// the object detection models take an input of [1,?,?,3]
func buildDecodeBitmapGraph() (g *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	channels := int64(3)

	input = op.Placeholder(s, tf.String)
	output = op.ExpandDims(s,
		op.DecodeJpeg(s, input, op.DecodeJpegChannels(channels)),
		op.Const(s.SubScope("make_batch"), int32(0)),
	)
	g, err = s.Finalize()
	return
}

func SavePng(path string, jp image.Image) error {
	fg, err := os.Create(path)
	defer fg.Close()
	if err != nil {
		return err
	}
	err = png.Encode(fg, jp)
	if err != nil {
		return err
	}
	return nil
}

func SaveJpeg(path string, jp image.Image) error {
	fg, err := os.Create(path)
	defer fg.Close()
	if err != nil {
		return err
	}
	err = jpeg.Encode(fg, jp, nil)
	if err != nil {
		return err
	}
	return nil
}

func AsGrayWithChannelWeights(img image.Image, r float64, g float64, b float64) image.Image {
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixelVal := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixelVal).(color.RGBA)

			// Offset colors a little, adjust it to your taste
			r := float64(originalColor.R) * 1.0
			g := float64(originalColor.G) * 1.0
			b := float64(originalColor.B) * 1.0

			// average
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}

	return wImg
}

func AsGray(img image.Image) image.Image {
	return AsGrayWithChannelWeights(img, 1.0, 1.0, 1.0)
}
