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

	// An Outcome simply is an element of a set,
	// namely, the OutcomeSpace
	Outcome set.Element

	// An Outcomes list is a slice of set.Element
	Outcomes []set.Element

	// An OutcomeSpace is a set
	OutcomeSpace set.Interface

	// A Distribution is the interface for interacting with
	// probability distribution. The domain is the universal set,
	// the outcomes are the support.
	Distribution interface {
		Domain() set.Interface
		Outcomes() set.Interface
		ProbabilityOf(Outcome) Probability
	}

	// A DiscreteDistribution is the interface for a distribution
	// we can programmatically manipulate. It inherits the interface
	// of a general ProbabilityDistribution, and adds the ability to add
	// outcomes.
	DiscreteDistribution interface {
		Distribution
		Support() Outcomes
		AddOutcome(Outcome, Probability)
	}

	// An Event is a set. As in probability theory, this set should
	// be a subset of the outcome space. e.g., an event A ⊆ Ω
	Event set.Interface

	// A RandomVariable is defined to be a real valued function of an outcome
	// of a random experiment.
	RandomVariable func(Outcome) float64
)

// --- }}}

// --- Probability (Modeling Uncertainty) {{{

const (
	// Impossible represents the probability that an outcome or event
	// will never occur.
	Impossible Probability = 0

	// Certain represents the probability that an outcome or event
	// will certainly occur.
	Certain Probability = 1
)

// Valid determines whether a Probability value is valid. All
// probabilities must be on the inverval: [0, 1]
func (p Probability) Valid() bool {
	return p >= 0 && p <= 1
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
// of the DiscreteDistribution (and therefor implicitly
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
	assert(Support(d) < 1, "distribution already fully supported")
	assert(Support(d)+p <= 1, "adding outcome would over-support")
	assert(p.Valid(), "invalid probability")

	d.outcomes.Add(o)
	d.support[o] = p
}

func (d *distribution) ProbabilityOf(o Outcome) Probability {
	p, ok := d.support[o]

	if ok {
		return p
	} else {
		return Impossible
	}
}

// --- }}}

// --- Distribution Properties {{{
func Support(d Distribution) Probability {
	p := Probability(0)

	for o := range d.Outcomes().Iter() {
		p += d.ProbabilityOf(o)
	}

	return p
}

func FullySupported(d Distribution) bool {
	return Support(d) == Certain
}

func Cardinality(d DiscreteDistribution) uint {
	return d.Outcomes().Cardinality()
}

func Degenerate(d DiscreteDistribution) bool {
	return Cardinality(d) == 1 && FullySupported(d)
}

// --- }}}

// --- Events {{{

func ProbabilityOf(d Distribution, A Event) Probability {
	sum := Impossible

	for a := range A.Iter() {
		sum += d.ProbabilityOf(a)
	}

	return sum
}

func IndependentEvents(d Distribution, A, B Event) bool {
	return ProbabilityOf(d, set.Union(A, B)) == ProbabilityOf(d, A)*ProbabilityOf(d, B)
}

// --- }}}

// --- Random Variables {{{

func Moment(d Distribution, X RandomVariable, n int) float64 {
	moment := func(o Outcome) float64 {
		return math.Pow(X(o), float64(n))
	}

	return Expectation(d, moment)
}

func Expectation(d Distribution, X RandomVariable) float64 {
	exp := 0.0

	for o := range d.Outcomes().Iter() {
		exp += X(o) * float64(d.ProbabilityOf(o))
	}

	return exp
}

func Variance(d Distribution, X RandomVariable) float64 {
	return Moment(d, X, 2) - math.Pow(Moment(d, X, 1), 2.0)
}

func Covariance(d Distribution, X, Y RandomVariable) float64 {
	return Expectation(d, func(o Outcome) float64 { return X(o) * Y(o) }) - Expectation(d, X)*Expectation(d, Y)
}

func IndependentVariables(d Distribution, X, Y RandomVariable) bool {
	return Covariance(d, X, Y) == 0
}

// --- }}}

// --- Composition {{{

func Compose(p, q DiscreteDistribution, alpha Probability) DiscreteDistribution {
	assert(FullySupported(p), "first distribution is not fully supported")
	assert(FullySupported(q), "second distribution is not fully supported")
	assert(set.Equivalent(p.Domain(), q.Domain()), "distribution domains must be equivalent")

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
