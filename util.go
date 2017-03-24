package main

import (
	"fmt"
	"image"
	"math"
)

// Convert a slice of Interfaces to a slice of Points.
func castImage(interfaces []interface{}) Image {
	var path = make(Image, len(interfaces))
	for i, v := range interfaces {
		path[i] = v.(Triangle)
	}
	return path
}

// Convert a slice of Points to a slice of interfaces.
func uncastImage(image Image) []interface{} {
	var interfaces = make([]interface{}, len(image))
	for i, v := range image {
		interfaces[i] = v
	}
	return interfaces
}

func FastCompare(img1, img2 *image.RGBA) (float64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := float64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += sqDiffUInt8(img1.Pix[i], img2.Pix[i])
	}

	return math.Sqrt(accumError), nil
}

func sqDiffUInt8(x, y uint8) float64 {
	d := float64(x) - float64(y)
	return d * d
}
