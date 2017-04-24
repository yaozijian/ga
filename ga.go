package main

import (
	"math/rand"
	"time"
)

func main() {

	var (
		pops Population
		ga   *GeneAlg
	)

	rand.Seed(time.Now().UnixNano())

	pops.init(50)

	ga = &GeneAlg{
		pplt:   pops,
		maxgas: 10000,
		crate:  0.9,
		mrate:  0.1,
	}

	ga.evolvePopulation()
}
