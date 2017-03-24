package main

import (
	"math/rand"
)

func (p Image) MutatePoint(rng *rand.Rand) {

	i := rng.Intn(len(p))
	j := rng.Intn(2)
	switch j {
	case 0:
		p[i].p1.x = rng.Float64()
		p[i].p1.y = rng.Float64()
	case 1:
		p[i].p2.x = rng.Float64()
		p[i].p2.y = rng.Float64()
	case 2:
		p[i].p3.x = rng.Float64()
		p[i].p3.y = rng.Float64()
	}
}

func (p Image) MutateColorOne(rng *rand.Rand) {
	i := rng.Intn(len(p))
	j := rng.Intn(3)

	switch j {
	case 0:
		p[i].c.R = uint8(rng.Intn(0xff))
	case 1:
		p[i].c.G = uint8(rng.Intn(0xff))
	case 2:
		p[i].c.B = uint8(rng.Intn(0xff))
	case 3:
		p[i].c.A = uint8(rng.Intn(0xff))
	}

}

func (p Image) MutSplice(rng *rand.Rand) {
	var split = rng.Intn(len(p)-1) + 1
	copy(p, append(p[split:], p[:split]...))
}

func (p Image) MutPermute(rng *rand.Rand) {

	var (
		n = 2
	)
	// Nothing to permute
	if len(p) <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		// Choose two points on the genome
		var (
			i = rng.Intn(len(p))
			j = rng.Intn(len(p))
		)
		// Permute the genes
		p[i], p[j] = p[j], p[i]
	}
}
