package main

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"
)

var (
	uninitialized float64 = math.Inf(1)
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
	total := uninitialized
	var operator Gene
	lookingForNumber := true

	for index := 0; index < len(c.body)*2; index++ {
		curr := c.GetGene(index)
		// fmt.Printf("curr %s\n", curr)
		if lookingForNumber && curr.IsNumber() {
			// fmt.Printf("handle %f\n", curr.GetFloatValue())
			lookingForNumber = false
			if total == uninitialized {
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

func (c *Chromosome) Mutate(rate float64) *Chromosome {
	for i := 0; i < len(c.body); i++ {
		for j := 0; j < 8; j++ {
			if 0 == mathrand.Int31n(int32(1/rate)) {
				bitmask := byte(1 << uint(j+1))
				if (c.body[i] & bitmask) > 0 {
					c.body[i] &= ^bitmask
				} else {
					c.body[i] |= bitmask
				}
			}
		}
	}
	return c
}
