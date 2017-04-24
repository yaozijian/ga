package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type (
	Chromosome interface { //染色体
		CrossWith(Chromosome) // 交叉
		Mutation()            // 变异
		Clone() Chromosome
	}

	Fitness interface { // 适应度
		BetterThan(Fitness) bool // 适应度比较
		Weight() float64         // 随机选择(轮盘赌)时的权重
		Clone() Fitness
	}

	Population []Chromosome // 种群

	GeneAlgParam struct {
		Pop            Population               // 种群
		CrossRate      float64                  // 交叉率[0,1.0)
		MutationRate   float64                  // 变异率[0,1.0)
		MaxGenerations int                      // 最大进化代数
		CalcFitness    func(Chromosome) Fitness // 计算适应度函数
	}

	GeneAlgResult struct {
		MaxFit Fitness    // 最大适应度
		MaxInd Chromosome // 具有最大适应度的个体
	}

	GeneAlgo interface {
		Run(GeneAlgParam) GeneAlgResult
	}

	//------------------------------------------------

	geneAlgo struct {
		r *rand.Rand
	}

	calcFitAndWeightContext struct {
		bestFit Fitness              // 最佳适应度
		bestInd Chromosome           // 最佳个体
		ind2ctx []*individualContext // 每个个体的上下文
	}

	individualContext struct {
		fitness Fitness // 适应度
		weight  float64 // 权重
		prop    float64 // 份额
	}
)

func NewGeneAlgo() GeneAlgo {
	var v int64
	binary.Read(crand.Reader, binary.LittleEndian, &v)
	x := &geneAlgo{r: rand.New(rand.NewSource(v))}
	return x
}

func (ga *geneAlgo) Run(param GeneAlgParam) GeneAlgResult {

	var (
		bestFit Fitness
		bestInd Chromosome
	)

	ctx := &calcFitAndWeightContext{}

	// 计算适应度和权重,以及份额
	ctx.calcFitnessAndWeight(param)

	for x := 0; x < param.MaxGenerations; x++ {

		// 记录最佳个体和最佳适应度
		if bestFit == nil || ctx.bestFit.BetterThan(bestFit) {
			bestFit = ctx.bestFit.Clone()
			bestInd = ctx.bestInd.Clone()
		}

		// 进化
		ga.nextGeneration(param, ctx)
	}

	return GeneAlgResult{
		MaxFit: bestFit,
		MaxInd: bestInd,
	}
}

func (ga *geneAlgo) nextGeneration(param GeneAlgParam, ctx *calcFitAndWeightContext) {

	// 随机选择两个个体进行交叉操作
	if ga.r.Float64() <= param.CrossRate {

		chrom1 := ctx.selection(ga.r)
		chrom2 := ctx.selection(ga.r)

		param.Pop[chrom1].CrossWith(param.Pop[chrom2])

		// 交叉操作会改变个体,进而改变群体,所以需要重新计算适应度和权重,以及份额
		ctx.calcFitnessAndWeight(param)
	}

	// 随机选择一个个体进行变异操作
	if ga.r.Float64() <= param.MutationRate {
		chrom3 := ctx.selection(ga.r)
		param.Pop[chrom3].Mutation()
		// 变异操作会改变个体,进而改变群体,所以需要重新计算适应度和权重,以及份额
		ctx.calcFitnessAndWeight(param)
	}
}

func (ctx *calcFitAndWeightContext) calcFitnessAndWeight(param GeneAlgParam) {

	var sum float64

	ctx.ind2ctx = make([]*individualContext, len(param.Pop)) // 种群序号 --> 适应度和权重
	ctx.bestFit = nil
	ctx.bestInd = nil

	// 计算适应度和权重
	for x := range param.Pop {

		fitness := param.CalcFitness(param.Pop[x])

		ctx.ind2ctx[x] = &individualContext{
			fitness: fitness,
			weight:  fitness.Weight(),
		}

		sum += ctx.ind2ctx[x].weight

		if ctx.bestFit == nil || fitness.BetterThan(ctx.bestFit) {
			ctx.bestFit = fitness.Clone()
			ctx.bestInd = param.Pop[x].Clone()
		}
	}

	// 每个个体所占份额
	var prop float64

	for _, item := range ctx.ind2ctx {
		prop += item.weight / sum
		item.prop = prop
	}
}

func (ctx *calcFitAndWeightContext) selection(r *rand.Rand) (selidx int) {

	// 随机(轮盘赌)确定选择哪个个体
	dice := r.Float64()

	size := len(ctx.ind2ctx)

	for i := 0; i < size; i++ {
		if dice < ctx.ind2ctx[i].prop {
			selidx = i
			break
		}
	}

	return
}
