package main

import (
	"math"
	"math/rand"
)

type (
	chromosomeImpl float64
	fitnessImpl    float64
)

func (c *chromosomeImpl) CrossWith(o Chromosome) {

	x := o.(*chromosomeImpl)

	m := (*c + *x) / 2.0

	if rand.Intn(2) == 0 {
		*c = m
	} else {
		*x = m
	}
}

func (c *chromosomeImpl) Mutation() {
	if rand.Intn(2) == 0 {
		*c += 0.01
	} else {
		*c -= 0.01
	}
}

func (c *chromosomeImpl) Clone() Chromosome {
	t := chromosomeImpl(*c)
	return &t
}

//-----------------------------------------------

func (f *fitnessImpl) BetterThan(o Fitness) bool {
	x := o.(*fitnessImpl)
	return *f > *x
}

func (f *fitnessImpl) Weight() float64 {
	return float64(*f)
}

func (f *fitnessImpl) Clone() Fitness {
	x := fitnessImpl(*f)
	return &x
}

//-----------------------------------------------

func calcFitness(c Chromosome) Fitness {
	x := float64(*(c.(*chromosomeImpl)))
	y := fitnessImpl(2 + x*math.Sin(10*math.Pi*x))
	return &y
}
