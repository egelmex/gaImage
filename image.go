package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Point struct {
	x, y float64
}

type Triangle struct {
	c          color.NRGBA
	p1, p2, p3 Point
}

type Image []Triangle

func randomColor(rng *rand.Rand) color.NRGBA {
	return color.NRGBA{
		uint8(rng.Intn(0xff)),
		uint8(rng.Intn(0xff)),
		uint8(rng.Intn(0xff)),
		uint8(rng.Intn(0xff))}
}

func MakeTriangle(rng *rand.Rand) Triangle {
	c := randomColor(rng)
	p1 := Point{rng.Float64(), rng.Float64()}
	p2 := Point{rng.Float64(), rng.Float64()}
	p3 := Point{rng.Float64(), rng.Float64()}
	return Triangle{c, p1, p2, p3}
}
func MakeImage(rng *rand.Rand) Image {
	var (
		image = make(Image, triangleCount)
	)
	for i := 0; i < triangleCount; i++ {
		image[i] = MakeTriangle(rng)
	}
	return image
}

func (frame Image) drawFrame() *image.RGBA {
	var (
		width, height = getTargetDimentions()
		w, h          = float64(width), float64(height)
	)
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetLineWidth(0)

	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.MoveTo(0.0*w, 0.0*h)
	gc.LineTo(0.0*w, 1.0*h)
	gc.LineTo(1.0*w, 1.0*h)
	gc.LineTo(1.0*w, 0.0*h)
	gc.Close()
	gc.FillStroke()

	for _, t := range frame {
		//r,g,b,a := t.c.RGBA()
		//gc.SetFillColor(color.RGBA{uint8(r),uint8(g),uint8(b),uint8(a)})
		gc.SetFillColor(t.c)
		gc.MoveTo(t.p1.x*w, t.p1.y*h)
		gc.LineTo(t.p2.x*w, t.p2.y*h)
		gc.LineTo(t.p3.x*w, t.p3.y*h)
		gc.Close()
		gc.FillStroke()

	}
	return dest
}
