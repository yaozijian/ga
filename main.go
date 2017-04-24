package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

func main() {

	pop := initPopulation(300)

	param := GeneAlgParam{
		Pop:            pop,
		CrossRate:      0.9,
		MutationRate:   0.1,
		MaxGenerations: 10000,
		CalcFitness:    calcFitness,
	}

	result := NewGeneAlgo().Run(param)

	fmt.Printf(
		"best solution: %f,best chromosome: %f\n",
		float64(*(result.MaxFit.(*fitnessImpl))),
		float64(*(result.MaxInd.(*chromosomeImpl))),
	)
}

func initPopulation(size int) []Chromosome {

	var v int64
	binary.Read(crand.Reader, binary.LittleEndian, &v)
	rand.Seed(v)

	pops := make([]Chromosome, size)

	for i := 0; i < size; i++ {
		x := new(chromosomeImpl)
		*x = chromosomeImpl(-1.0 + 3*rand.Float64())
		pops[i] = x
	}

	return pops
}
