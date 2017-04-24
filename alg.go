package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type (
	Chromosome struct {
		gene    float64
		fitness float64
		rank    float64
	}

	Population []*Chromosome

	GeneAlg struct {
		pplt   Population
		maxgas int
		crate  float64
		mrate  float64
	}
)

func fitness(x float64) float64 {
	return 2 + x*math.Sin(10*math.Pi*x)
}

func (p *Population) init(size int) {

	pops := make([]*Chromosome, size)

	for i := 0; i < size; i++ {
		pops[i] = new(Chromosome)
		pops[i].SetGene(-1.0 + 3*rand.Float64())
	}

	*p = Population(pops)
}

func (c *Chromosome) SetGene(gene float64) {
	c.gene = gene
	c.fitness = fitness(gene)
}

func (c *Chromosome) Fitness() float64 {
	return c.fitness
}

func (ga *GeneAlg) evolvePopulation() {

	var (
		counter         int
		best            float64
		globalBest      float64
		chromBest       Chromosome
		globalChromBest Chromosome
	)

	for counter++; counter <= ga.maxgas; counter++ {

		for i := 0; i < len(ga.pplt); i++ {
			cur := ga.pplt[i]
			cur.SetGene(cur.gene)
			if cur.Fitness() > best {
				best = cur.Fitness()
				chromBest = *cur
			}
		}

		if best > globalBest {
			globalBest = best
			globalChromBest = chromBest
		}

		ga.nextGeneration()
	}

	fmt.Printf("best solution: %.3f,best chromosome: %.3f\n", globalChromBest.fitness, globalChromBest.gene)
}

func (ga *GeneAlg) nextGeneration() {
	chrom1 := ga.selection()
	chrom2 := ga.selection()
	ga.crossover(chrom1, chrom2)
	ga.mutation(ga.selection())
}

func (ga *GeneAlg) selection() int {

	sort.SliceStable(ga.pplt, func(x, y int) bool {
		return ga.pplt[x].Fitness() < ga.pplt[y].Fitness()
	})

	var (
		sum   float64
		chrom int
	)

	size := len(ga.pplt)

	// 总份额
	for i := 0; i < size; i++ {
		sum += ga.pplt[i].fitness
	}

	// 每个个体所占份额
	for i := 0; i < size; i++ {
		ga.pplt[i].rank = ga.pplt[i].fitness / sum
	}

	// 随机(轮盘赌)确定选择哪个个体
	r := rand.Float64()

	for i := 0; i < size-1; i++ {
		if ga.pplt[i].rank <= r && r < ga.pplt[i+1].rank {
			chrom = i
			break
		}
	}

	return chrom
}

// 交叉
func (ga *GeneAlg) crossover(index1, index2 int) {
	if rand.Float64() <= ga.crate {
		gene := (ga.pplt[index1].gene + ga.pplt[index2].gene) / 2
		if rand.Intn(2) == 0 {
			ga.pplt[index1].SetGene(gene)
		} else {
			ga.pplt[index2].SetGene(gene)
		}
	}
}

// 变异
func (ga *GeneAlg) mutation(index int) {
	if rand.Float64() <= ga.mrate {
		if rand.Intn(2) == 0 {
			ga.pplt[index].gene += 0.01
		} else {
			ga.pplt[index].gene -= 0.01
		}
	}
}
