package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/nlandolfi/prob"
	"github.com/nlandolfi/set"
)

func main() {
	log.Printf("Factorial(80) = %s", prob.Factorial(big.NewInt(80)))
	log.Printf("Factorial(4) = %s", prob.Factorial(big.NewInt(4)))
	log.Printf("Factorial(76) = %s", prob.Factorial(big.NewInt(76)))
	log.Printf("Combination(80, 4) = %s", prob.Combination(big.NewInt(80), big.NewInt(4)))
	/*
		log.Printf(prob.Combination(50, 3))
		log.Printf(prob.Combination(30, 1))
		log.Printf(prob.Combination(50, 3) * prob.Combination(30, 1) / prob.Combination(80, 4))
	*/

	diceOutcomes := set.With([]set.Element{1, 2, 3, 4, 5, 6})

	d := prob.NewUniformDiscrete(diceOutcomes)

	var x prob.RandomVariable = func(o prob.Outcome) float64 {
		return float64(o.(int))
	}

	log.Printf("Expectation of dice roll: %f", prob.Expectation(d, x))
	log.Printf("Variance of dice roll: %f", prob.Variance(d, x))

	u4 := prob.NewUniformDiscrete(set.WithElements(1, 2, 3, 4))

	X := func(o prob.Outcome) float64 {
		switch int(o.(int)) {
		case 1:
			return 1
			break
		case 2:
			return 0
			break
		case 3:
			return -1
			break
		case 4:
			return 0
			break
		default:
			panic(fmt.Sprintf("invalid outcome %+v", o))
		}
		return 0
	}

	Y := func(o prob.Outcome) float64 {
		switch int(o.(int)) {
		case 1:
			return 0
			break
		case 2:
			return 1
			break
		case 3:
			return 0
			break
		case 4:
			return -1
			break
		default:
			panic(fmt.Sprintf("invalid outcome %+v", o))
		}

		return 0
	}

	log.Printf("Are X and Y independent? (E[XY] - E[X]E[Y] ?=? 0), %t", prob.Independent(u4, X, Y))

	e1 := prob.Binomial(2, 0.5)(0) == prob.Geometric(.5)(2)
	e2 := prob.Geometric(0.5)(2) == prob.Uniform(4)(1)

	log.Printf("B(2, 0.5)(0) = %f", prob.Binomial(2, 0.5)(0))
	log.Printf("G(0.5)(2) = %f", prob.Geometric(0.5)(2))
	log.Printf("U(4)(1) = %f", prob.Uniform(4)(1))
	log.Printf("B(2, 0.5)(0) == G(.5)(2) == U(4)(1), %t", e1 && e2)
}
