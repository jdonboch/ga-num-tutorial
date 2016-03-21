package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
)

type Chromosome struct {
	body []byte
}

func GenerateRandomChromosome(numberOfGenes int) *Chromosome {
	chromo := &Chromosome{
		body: make([]byte, numberOfGenes/2),
	}
	rand.Read(chromo.body)
	return chromo
}

func (c *Chromosome) GetGene(index int) Gene {
	b := c.body[index/2]

	if (index % 2) == 0 {
		return Gene(b >> 4)
	}
	return Gene(b & 0x0f)
}

func (c *Chromosome) CalculateTotal() float64 {
	const uninit float64 = float64(-77777)
	total := uninit
	var operator Gene
	lookingForNumber := true

	for index := 0; index < len(c.body)*2; index++ {
		curr := c.GetGene(index)
		// fmt.Printf("curr %s\n", curr)
		if lookingForNumber && curr.IsNumber() {
			// fmt.Printf("handle %f\n", curr.GetFloatValue())
			lookingForNumber = false
			if total == uninit {
				total = curr.GetFloatValue()
			} else {
				total = operator.Operate(total, curr.GetFloatValue())
			}
		} else if !lookingForNumber && curr.IsOperator() {
			// fmt.Printf("handle %s\n", curr)
			lookingForNumber = true
			operator = curr
		}
	}
	return total
}

func (c *Chromosome) GetFitnessScore(target float64) float64 {
	total := c.CalculateTotal()
	if math.IsNaN(total) || math.IsInf(total, 0) {
		return 0
	}
	return math.Abs(1 / (target - total))
}

func (c *Chromosome) Mate(other *Chromosome, crossoverRate float64) *Chromosome {
	if mathrand.Float64() <= crossoverRate {
		length := len(c.body)
		return &Chromosome{
			body: append(c.body[:length/2], other.body[length/2:]...),
		}
	}
	return c
}

func (c *Chromosome) Mutate(rate float64) {

}

func SumPopulationProbability(pop []*Chromosome, target float64) ([]float64, int) {
	sumOfFitness := big.NewFloat(0.0)
	var sumOfProb float64
	for index, chromo := range pop {
		score := chromo.GetFitnessScore(target)
		if math.IsInf(score, 0) || math.IsNaN(score) {
			fmt.Printf("Inf\n")
			return nil, index
		}
		sumOfFitness.Add(sumOfFitness, big.NewFloat(score))
	}

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

	// fmt.Printf("Prob slice: %+v\n", probSlice)

	return probSlice, -1
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

const targetSolution = 145

const numGenes int = 10
const initalPopulateSize int = 1000
const finalPopulation int = 10000
const crossoverRate float64 = 0.7
const mutationRate float64 = 0.001

func main() {

	population := make([]*Chromosome, initalPopulateSize)
	for i := 0; i < initalPopulateSize; i++ {
		population[i] = GenerateRandomChromosome(numGenes)
	}

	for len(population) < finalPopulation {
		var mates [2]*Chromosome

		probSlice, solution := SumPopulationProbability(population, targetSolution)
		if solution >= 0 {
			OutputSolution(population[solution])
			break
		}
		for j := 0; j < 2; j++ {
			chosenProb := mathrand.Float64()
			// fmt.Printf("Chosen prob: %f\n", chosenProb)
			for k, prob := range probSlice {
				if chosenProb <= prob {
					mates[j] = population[k]
					break
				}
			}
		}

		population = append(population, mates[0].Mate(mates[1], crossoverRate))
		// TODO Mutate
	}
	OutputSolution(BestSolution(population, targetSolution))
}

func OutputSolution(c *Chromosome) {
	out := ""
	for i := 0; i < numGenes; i++ {
		out += fmt.Sprintf("%s ", c.GetGene(i))
	}
	fmt.Println(out)
	fmt.Println(c.CalculateTotal())
	fmt.Println(c.GetFitnessScore(targetSolution))
}
