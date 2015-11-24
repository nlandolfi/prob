package prob

import (
	"math"
	"math/big"
)

func Binomial(n int64, p Probability) func(int64) Probability {
	return func(k int64) Probability {
		return Probability(float64(nint(0).Binomial(n, k).Int64()) * math.Pow(float64(p), float64(k)) * math.Pow(1-float64(p), float64(n-k)))
	}
}

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

func Uniform(n int) func(int) Probability {
	return func(k int) Probability {
		return Probability(1.0 / float64(n))
	}
}

func Geometric(p Probability) func(int) Probability {
	return func(k int) Probability {
		return Probability(math.Pow(float64(Certain-p), float64(k-1)) * float64(p))
	}
}

func Poisson(mu float64) func(int) Probability {
	return func(k int) Probability {
		return Probability(math.Pow(math.E, -mu) * math.Pow(mu, float64(k)) / float64(Factorial(big.NewInt(int64(k))).Int64()))
	}
}

func nint(i int64) *big.Int {
	return big.NewInt(i)
}

func Factorial(n *big.Int) *big.Int {
	if n.Cmp(nint(0)) == 0 {
		return big.NewInt(1)
	}

	i, z := nint(0), nint(0)

	return z.Mul(n, Factorial(i.Sub(n, nint(1))))
}

func Combination(n, k *big.Int) *big.Int {
	delta, z := nint(0), nint(0)
	delta.Sub(n, k)
	return z.Div(Factorial(n), z.Mul(Factorial(k), Factorial(delta)))
}

func Choose(n, k *big.Int) *big.Int {
	return Combination(n, k)
}
