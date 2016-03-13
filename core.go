package prob

import (
	"math"
	"math/rand"

	"github.com/nlandolfi/set"
)

// --- Types {{{
type (
	// A Probability is an element of [0, 1]
	Probability float64

	// An Outcome is an element of a set,
	// namely the Outcome Space
	Outcome set.Element

	// An Outcomes list is a slice of set.Element
	Outcomes []set.Element

	// An OutcomeSpace is a set
	OutcomeSpace set.Interface

	// A Distribution is the interface for interacting with
	// probability distribution. The domain is the universal set,
	// for that distribution, the outcomes are the support.
	Distribution interface {
		// Domain is the set which defines the possible outcomes of
		// a Distribution. This is the Outcome Space
		Domain() set.Interface

		// Outcomes is the set of Outcomes in the Domain which occur
		// with a non-zero probability
		Outcomes() set.Interface

		// ProbabilityOf returns the probability of a given Outcome
		//
		// Note: ProbabilityOf should return a probability of 0 (Impossible)
		// for any outcome in the domain, but without defined support
		ProbabilityOf(Outcome) Probability
	}

	// A DiscreteDistribution is the interface for a distribution
	// we can programmatically manipulate. It inherits the interface
	// of a general Distribution
	//
	// Note: We can examine/add outcomes
	DiscreteDistribution interface {
		Distribution
		Support() Outcomes
		AddOutcome(Outcome, Probability)
	}

	// An Event is a set. As in probability theory, this set should
	// be a subset of the outcome space. e.g., an event A ⊆ Ω
	Event set.Interface

	// A RandomVariable is defined to be a real valued function of an outcome
	// of a random experiment. It is neither random, nor variable. It is a
	// fixed function. The Outcome introduces the stochasticity.
	RandomVariable func(Outcome) float64
)

// --- }}}

// --- Probability {{{

const (
	// Impossible represents the probability that an outcome or event
	// will never occur.
	Impossible Probability = 0.0

	// Certain represents the probability that an outcome or event
	// will certainly occur.
	Certain Probability = 1.0
)

// Valid determines whether a Probability value is valid. All
// probabilities must be on the inverval: [0, 1]
func (p Probability) Valid() bool {
	return p >= 0 && p <= 1
}

// epsilon is the acceptable floating point error
var epsilon = 0.00001

// equiv determines whether two float64s are equivalent to each
// other with respect to epsilon
func equiv(f1, f2 float64) bool {
	return math.Abs(f1-f2) < epsilon
}

// --- }}}

// --- Discrete Distribution --- {{{

// NewDiscreteDistribution constructs a discrete distribution over
// the set d, provided as the domain of the distribution
func NewDiscreteDistribution(d set.Interface) DiscreteDistribution {
	return &distribution{
		domain:   d,
		outcomes: set.New(),
		support:  make(map[Outcome]Probability),
	}
}

// NewUniformDiscrete constructs a discrete distribution over the
// set domain. Each element is assigned a probability 1/Cardinality(domain)
func NewUniformDiscrete(domain set.Interface) DiscreteDistribution {
	d := NewDiscreteDistribution(domain)

	individualSupport := Certain / Probability(domain.Cardinality())

	for o := range domain.Iter() {
		d.AddOutcome(o, individualSupport)
	}

	return d
}

// distribution structure serves as an implementation
// of the DiscreteDistribution (and therefore implicitly
// Distribution) interfaces
type distribution struct {
	domain   set.Interface
	outcomes set.Interface
	support  map[Outcome]Probability
}

func (d *distribution) Domain() set.Interface {
	return d.domain
}

func (d *distribution) Outcomes() set.Interface {
	return d.outcomes
}

func (d *distribution) Support() Outcomes {
	return d.outcomes.Elements()
}

func (d *distribution) AddOutcome(o Outcome, p Probability) {
	assert(!equiv(float64(Support(d)), 1.0), "distribution already fully supported")
	assert((float64(Support(d)+p) < 1.0+epsilon), "adding outcome would over-support")
	assert(p.Valid(), "invalid probability")
	assert(!equiv(float64(p), 0), "probability zero")

	d.outcomes.Add(o)
	d.support[o] = p
}

func (d *distribution) ProbabilityOf(o Outcome) Probability {
	p, ok := d.support[o]

	if ok {
		return p
	}

	if d.domain.Contains(o) {
		return Impossible
	} else {
		panic("outcome not in domain")
	}
}

// --- }}}

// --- Distribution Properties {{{

// Support calculates the total portion of the probability mass
// we have defined. Recall that all distributions have a mass of
// 1.
//
// I.e., if the Support(d) = 1, d ∈ Distributions, then we say
// that d is _fully-supported_. Adding another outcome with
// a non-zero probability would invalidate the distribution.
func Support(d Distribution) Probability {
	p := Probability(0.0)

	for o := range d.Outcomes().Iter() {
		p += d.ProbabilityOf(o)
	}

	return p
}

// FullySupported checks that a Distribution has assigned all
// of it's probability mass.
//
// True iff the sum over the probabilities of all outcomes is 1.
func FullySupported(d Distribution) bool {
	return equiv(float64(Support(d)), float64(Certain))
}

// The cardinality of a discrete distribution is the number of
// potential outcomes
func Cardinality(d DiscreteDistribution) uint {
	return d.Outcomes().Cardinality()
}

// Degenerate evaluates whether a DiscreteDistribution is degenerate.
//
// A degenerate distribution is fully supported, but with only one
// outcome.
func Degenerate(d DiscreteDistribution) bool {
	return Cardinality(d) == 1 && FullySupported(d)
}

// --- }}}

// --- Events {{{

// ProbabilityOf calculates the probability of an event, A, given
// a Distribution, d
func ProbabilityOf(d Distribution, A Event) Probability {
	sum := Impossible

	for a := range A.Iter() {
		sum += d.ProbabilityOf(a)
	}

	return sum
}

// IndependentEvents determines whether A and B are independent
// under the distribution d.
//
// Equivalently: P(A, B) = P(A)P(B)
func IndependentEvents(d Distribution, A, B Event) bool {
	return equiv(float64(ProbabilityOf(d, set.Union(A, B))), float64(ProbabilityOf(d, A)*ProbabilityOf(d, B)))
}

// --- }}}

// --- Random Variables {{{

// Moment calculates the nth moment of a random variable
//
// Recally: the nth moment of a random variable X over a
// distribution d is the expectation of X^n
func Moment(d Distribution, X RandomVariable, n int) float64 {
	moment := func(o Outcome) float64 {
		return math.Pow(X(o), float64(n))
	}

	return Expectation(d, moment)
}

// Expectation computes the expected value of a random variable,
// X over a distribution d
func Expectation(d Distribution, X RandomVariable) float64 {
	exp := 0.0

	for o := range d.Outcomes().Iter() {
		exp += X(o) * float64(d.ProbabilityOf(o))
	}

	return exp
}

// Variance computes the variance of a random variable, X,
// over a distribution d
//
// Recall: Var(X) = E(X^2) - E(X)^2
func Variance(d Distribution, X RandomVariable) float64 {
	return Moment(d, X, 2) - math.Pow(Moment(d, X, 1), 2.0)
}

// Covariance computes the covariance of the random variables X and Y,
// over a distribution d.
//
// Recall: Cov(X, Y) = E(XY) - E(X)E(Y)
func Covariance(d Distribution, X, Y RandomVariable) float64 {
	return Expectation(d, func(o Outcome) float64 { return X(o) * Y(o) }) - Expectation(d, X)*Expectation(d, Y)
}

// IndependentVariables determines whether two random variables X and Y are
// independent over the distribution d.
//
// Recall: X ind. Y iff Cov(X, Y) = 0
func IndependentVariables(d Distribution, X, Y RandomVariable) bool {
	return Covariance(d, X, Y) == 0
}

// --- }}}

// --- Composition {{{

// Compose takes two distributions, p and q, and creates a third lottery,
// n which takes on each outcome in p with a probability alpha times that
// events previous probability and each event in q with a probability 1-alpha
// times that events previous probability.
//
// for o ∈ p.Domain() intersect q.Domain(); P(x in n) is alpha*P(x in p) + (1-alpha)P(x in q)
func Compose(p, q DiscreteDistribution, alpha Probability) DiscreteDistribution {
	assert(FullySupported(p), "first distribution is not fully supported")
	assert(FullySupported(q), "second distribution is not fully supported")
	assert(set.Equivalent(p.Domain(), q.Domain()), "domains of both distributions must be equivalent")

	n := NewDiscreteDistribution(p.Domain())

	for o := range n.Domain().Iter() {
		cp := alpha*p.ProbabilityOf(o) + (1-alpha)*q.ProbabilityOf(o)
		if cp == Impossible {
			continue // don't bother supporting
		}

		n.AddOutcome(o, cp)
	}

	return n
}

// --- }}}

// --- Simulation {{{

// Simulate simulates an experiment with the distribution
// defined by the DiscreteDistribution
//
//		s := set.WithElements(1, 2, 3)
//		d := NewUniformDiscrete(s)
//		Simulate(d) => 1 w.p. 1/3, 2 w.p. 1/3, 3 w.p. 1/3
func Simulate(d DiscreteDistribution) Outcome {
	assert(FullySupported(d), "discrete distribution not fully supported")

	f := Probability(rand.Float64())
	p := Probability(0)

	var last Outcome
	for o := range d.Outcomes().Iter() {
		p += d.ProbabilityOf(o)
		last = o

		if f < p {
			return o
		}
	}

	return last
}

// --- }}}
