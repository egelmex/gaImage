package main

import (
	"flag"
	"fmt"
	"github.com/Maxhalford/gago"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	triangleCount = 50
	generations   = 50000
	//target           = MakeSample()
	target *image.RGBA
	//target           = loadImage("/home/me92/projects/gopath/src.jpg")
	printGenerations = 500
)

func getTargetDimentions() (int, int) {
	return target.Bounds().Dx(), target.Bounds().Dy()
}

type Point struct {
	x, y float64
}

type Triangle struct {
	c          color.NRGBA
	p1, p2, p3 Point
}

type Image []Triangle

func drawFrame(frame Image) *image.RGBA {
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

func (p Image) Evaluate() (distance float64) {
	frame := drawFrame(p)
	d, _ := FastCompare(frame, target)
	if math.IsInf(d, 0) {
		panic("inf!")
	}
	distance = d

	if d == 0 {
		panic("no way")
	}
	return
}

func (p Image) Mutate(rng *rand.Rand) {
	type mutFunc func(*rand.Rand)

	mutations := []mutFunc{
		p.MutatePoint,
		p.MutateColorOne,
		p.MutSplice,
		p.MutPermute,
	}

	n := rand.Intn(len(mutations))
	mutations[n](rng)

}

// Crossover a Path with another Path by using Partially Mixed Crossover (PMX).
func (p Image) Crossover(p1 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNX(uncastImage(p), uncastImage(p1.(Image)), 2, rng)
	return castImage(o1), castImage(o2)
}

func loadImage(path string) *image.RGBA {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	i, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	result := image.NewRGBA(i.Bounds())
	draw.Draw(result, i.Bounds(), i, image.ZP, draw.Src)

	return result
}

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

func MakeImage(rng *rand.Rand) gago.Genome {
	var (
		image = make(Image, triangleCount)
	)
	for i := 0; i < triangleCount; i++ {
		image[i] = MakeTriangle(rng)
	}
	return image
}

func doFlags() string {
	var (
		svar string
	)
	flag.StringVar(&svar, "path", "src.jpg", "Path to input file. Must be 1:1")
	flag.Parse()
	return svar

}

func configureGA(indervidualCount int) gago.GA {
	return gago.GA{
		MakeGenome: MakeImage,
		Topology: gago.Topology{
			NPopulations: 1,
			NIndividuals: indervidualCount,
		},
		Model: gago.ModSteadyState{
			Selector: gago.SelTournament{
				NParticipants: 4,
			},
			KeepBest: true,
			MutRate:  0.5,
		},
	}
}

func main() {
	runtime.GOMAXPROCS(4)

	var (
		ga        = configureGA(50)
		targetSrc = doFlags()
		date      = time.Now()
		folder    = date.Format("2006-01-02T15:04:05-0700")
	)

	fmt.Printf("Loading sample from: %s", targetSrc)
	target = loadImage(targetSrc)

	err := os.Mkdir(folder, 0777)
	if err != nil {
		panic(err)
	}

	ga.Initialize()
	ga.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err = ga.Validate()
	if err != nil {
		fmt.Printf("%s", err)
	}

	for i := 1; i < generations; i++ {
		ga.Enhance()

		//for _, population := range ga.Populations {
		//		for _, individual := range population.Individuals {
		//			f := fmt.Sprintf("%s/%d", folder, i)
		//			err := os.Mkdir(f, 0777)
		//			img := individual.Genome.(Image)
		//			filename := fmt.Sprintf("%s/%f.png", f, individual.Fitness)
		//			err = draw2dimg.SaveToPngFile(filename, drawFrame(img))
		//			if err != nil {
		//				panic(err)
		//			}
		//		}
		//	}
		if printGenerations > 0 && (i%printGenerations) == 0 {
			img := ga.Best.Genome.(Image)
			filename := fmt.Sprintf("%s/gen.%d.%d.png", folder, i, int(ga.Best.Fitness))
			err = draw2dimg.SaveToPngFile(filename, drawFrame(img))
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
	// Concatenate the elements from the best individual and display the result
	img := ga.Best.Genome.(Image)
	draw2dimg.SaveToPngFile("target.png", target)
	draw2dimg.SaveToPngFile("hello.png", drawFrame(img))
}
