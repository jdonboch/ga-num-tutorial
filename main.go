package main

// TODO
// * Determine why no convergence on solution
// * Use profiler to help performance

import (
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"os"
	"sort"
	"time"
)

func SumFitnessScore(pop []*Chromosome, target float64) (*big.Float, int) {
	sumOfFitness := big.NewFloat(0.0)
	for index, chromo := range pop {
		score := chromo.GetFitnessScore(target)
		if math.IsInf(score, 0) || math.IsNaN(score) {
			fmt.Printf("Inf\n")
			return nil, index
		}
		sumOfFitness.Add(sumOfFitness, big.NewFloat(score))
	}
	return sumOfFitness, -1
}

func SumPopulationProbability(pop []*Chromosome, target float64, sumOfFitness *big.Float) ([]float64, int) {
	var sumOfProb float64

	probSlice := make([]float64, len(pop))
	for index, chromo := range pop {

		bigFitScore := big.NewFloat(chromo.GetFitnessScore(target))
		prob, _ := bigFitScore.Quo(bigFitScore, sumOfFitness).Float64()
		probSlice[index] = sumOfProb + prob
		if math.IsNaN(probSlice[index]) {
			fmt.Printf("NaN,%f,%f,%v\n", sumOfProb, chromo.GetFitnessScore(target), sumOfFitness)
			return nil, index
		}
		sumOfProb += probSlice[index]
	}

	return probSlice, -1
}

func RouletteWheelSelection(pop []*Chromosome, target float64, sumOfFitness *big.Float) [2]*Chromosome {
	var mates [2]*Chromosome

	probSlice, solution := SumPopulationProbability(pop, targetSolution, sumOfFitness)
	if solution >= 0 {
		OutputSolution(pop[solution])
		os.Exit(0)
	}
	for j := 0; j < 2; j++ {
		chosenProb := mathrand.Float64()

		index := sort.SearchFloat64s(probSlice, chosenProb)
		mates[j] = pop[index]
	}

	return mates
}

func BestSolution(pop []*Chromosome, target float64) *Chromosome {
	var bestScore float64
	var bestChromo *Chromosome
	for _, chromo := range pop {
		currScore := chromo.GetFitnessScore(target)
		if currScore > bestScore {
			bestChromo = chromo
			bestScore = currScore
		}
	}
	return bestChromo
}

const (
	targetSolution     float64 = 146.0
	numGenes           int     = 12
	initalPopulateSize int     = 200
	finalPopulation    int     = 5000
	crossoverRate      float64 = 0.8
	mutationRate       float64 = 0.001
)

func main() {

	fmt.Println("Generating random chromosomes")

	population := make([]*Chromosome, initalPopulateSize)
	for i := 0; i < initalPopulateSize; i++ {
		population[i] = GenerateRandomChromosome(numGenes)
	}

	fmt.Println("Finished random chromosome gen")
	fmt.Println("Best solution in initial population")
	OutputSolution(BestSolution(population, targetSolution))
	fmt.Println("------------")

	mathrand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	sumOfFitness, solution := SumFitnessScore(population, targetSolution)
	if solution >= 0 {
		OutputSolution(population[solution])
		os.Exit(0)
	}

	for len(population) < finalPopulation {
		mates := RouletteWheelSelection(population, targetSolution, sumOfFitness)
		newChromo := mates[0].Mate(mates[1], crossoverRate).Mutate(mutationRate)
		sumOfFitness.Add(sumOfFitness, big.NewFloat(newChromo.GetFitnessScore(targetSolution)))
		population = append(population, newChromo)
	}

	fmt.Printf("Duration: %v\n", time.Since(start))

	OutputSolution(BestSolution(population, targetSolution))
}

func OutputSolution(c *Chromosome) {
	out := ""
	for i := 0; i < numGenes; i++ {
		out += fmt.Sprintf("%s ", c.GetGene(i))
	}
	fmt.Println(out)
	fmt.Printf("Expected: %.2f, Best Solution: %.2f\n", targetSolution, c.CalculateTotal())
}
