package prob

import (
	"math"
	"math/big"
)

// Bernoulli represents a Bernoulli trial
// { 1 with probability p, 0 with probability 1 - p }
func Bernoulli(p Probability) func(k int) Probability {
	return func(k int) Probability {
		if k == 1 {
			return p
		}

		return 1 - p
	}
}

// A Binomial distribution. The number of successes in n independent trials
// with a probability, p, of success in each trial.
// (n choose k)(p)^(k)(1-p)^(n-k)
func Binomial(n int64, p Probability) func(int64) Probability {
	return func(k int64) Probability {
		return Probability(float64(nint(0).Binomial(n, k).Int64()) * math.Pow(float64(p), float64(k)) * math.Pow(1-float64(p), float64(n-k)))
	}
}

// A Multinomial distribution. The number of elements in each category
// where the probability of being in category i is probabilities[i].
func Multinomial(probabilities ...Probability) func(...int) Probability {
	return func(partition ...int) Probability {
		assert(len(probabilities) == len(partition), "invalid partition")

		sum := 0
		for i := range partition {
			sum += partition[i]
		}

		assert(sum != 0, "partition sum can't be zero")

		num := Factorial(nint(int64(sum)))
		den := nint(1)

		for i := range partition {
			den.Mul(den, Factorial(nint(int64(partition[i]))))
		}

		scale := 1.0
		for i := range probabilities {
			scale *= math.Pow(float64(probabilities[i]), float64(partition[i]))
		}

		return Probability(float64(num.Div(num, den).Int64()) * scale)
	}
}

// A Uniform distribution on the discrete range [1, 2, ..., n]
func Uniform(n int) func(int) Probability {
	return func(k int) Probability {
		return Probability(1.0 / float64(n))
	}
}

// A Geometric distribution with parameter p.
//
// Recall that the geometric distribution models the probability that
// it takes k trials until we observe a success, where probability of a
// success in p
func Geometric(p Probability) func(int) Probability {
	return func(k int) Probability {
		return Probability(math.Pow(float64(Certain-p), float64(k-1)) * float64(p))
	}
}

// A Poisson distribution with paramter mu.
//
// Recall that the poisson distribution models the probability that we
// observe k successes in infinite trials; In other words, it models
// the expected number of occurrences in an interval of time t of a randomly
// occuring process with rate mu per t.
func Poisson(mu float64) func(int) Probability {
	return func(k int) Probability {
		return Probability(math.Pow(math.E, -mu) * math.Pow(mu, float64(k)) / float64(Factorial(big.NewInt(int64(k))).Int64()))
	}
}

// nint is a helper for big.NewInt
func nint(i int64) *big.Int {
	return big.NewInt(i)
}

// Factorial computes n!
func Factorial(n *big.Int) *big.Int {
	if n.Cmp(nint(0)) == 0 {
		return big.NewInt(1)
	}

	i, z := nint(0), nint(0)

	return z.Mul(n, Factorial(i.Sub(n, nint(1))))
}

// Combintation comuptes (n choose k)
func Combination(n, k *big.Int) *big.Int {
	delta, z := nint(0), nint(0)
	delta.Sub(n, k)
	return z.Div(Factorial(n), z.Mul(Factorial(k), Factorial(delta)))
}

// Choose is an alias of Combination
func Choose(n, k *big.Int) *big.Int {
	return Combination(n, k)
}
