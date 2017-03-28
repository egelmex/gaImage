package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/Maxhalford/gago"
	"github.com/llgcode/draw2d/draw2dimg"
)

type config struct {
	target string
}

var (
	triangleCount = 50
	generations   = 50000
	//target           = MakeSample()
	target *image.RGBA
	//target           = loadImage("/home/me92/projects/gopath/src.jpg")
	printGenerations = 500
	mutations        []mutFunc
)

func getTargetDimentions() (int, int) {
	return target.Bounds().Dx(), target.Bounds().Dy()
}

func (p Image) Evaluate() (distance float64) {
	frame := p.drawFrame()
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
	n := rand.Intn(len(mutations))
	mutations[n](p, rng)
}

// Crossover a Path with another Path by using Partially Mixed Crossover (PMX).
func (p Image) Crossover(p1 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNX(uncastImage(p), uncastImage(p1.(Image)), 2, rng)
	return castImage(o1), castImage(o2)
}

func doFlags() string {
	var (
		svar string
	)
	flag.StringVar(&svar, "path", "src.jpg", "Path to input file. Must be 1:1")

	mutations = []mutFunc{
		Image.MutSplice,
		Image.MutPermute,
	}

	var mPoint = *flag.Bool("mPoint", false, "Enable point mutation")
	var mColorRGBA = *flag.Bool("mColorRGBA", false, "Enable single channel RBGA mutation")

	flag.Parse()

	if mPoint {
		mutations = append(mutations, Image.MutatePoint)
	}
	if mColorRGBA {
		mutations = append(mutations, Image.MutateColorOne)
	}
	return svar

}

func configureGA(indervidualCount int) gago.GA {
	return gago.GA{
		MakeGenome: func(rng *rand.Rand) gago.Genome {
			return MakeImage(rng)
		},
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

func printGeneration(ga gago.GA, generation int, folder string) {
	for _, population := range ga.Populations {
		for _, individual := range population.Individuals {
			f := fmt.Sprintf("%s/%d", folder, generation)
			err := os.Mkdir(f, 0777)
			img := individual.Genome.(Image)
			filename := fmt.Sprintf("%s/%f.png", f, individual.Fitness)
			err = draw2dimg.SaveToPngFile(filename, img.drawFrame())
			if err != nil {
				panic(err)
			}
		}
	}
}

func printBest(ga gago.GA, generation int, folder string) {

	img := ga.Best.Genome.(Image)
	filename := fmt.Sprintf("%s/gen.%d.%d.png", folder, generation, int(ga.Best.Fitness))
	err := draw2dimg.SaveToPngFile(filename, img.drawFrame())
	if err != nil {
		panic(err)
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
		//printGenerations(ga, i, folder)
		printBest(ga, i, folder)

		if printGenerations > 0 && (i%printGenerations) == 0 {
		}
	}
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
	// Concatenate the elements from the best individual and display the result
	img := ga.Best.Genome.(Image)
	draw2dimg.SaveToPngFile("target.png", target)
	draw2dimg.SaveToPngFile("hello.png", img.drawFrame())
}
