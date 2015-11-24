package prob

import (
	"math"

	"github.com/nlandolfi/set"
)

type (
	Probability float64

	Outcome set.Element

	Outcomes []set.Element

	OutcomeSpace set.Interface

	Distribution interface {
		Domain() set.Interface
		Outcomes() set.Interface
		ProbabilityOf(Outcome) Probability
	}

	DiscreteDistribution interface {
		Distribution
		Support() Outcomes
		AddOutcome(Outcome, Probability)
	}

	RandomVariable func(Outcome) float64
)

// --- Probability (Modeling Uncertainty) {{{

const (
	Impossible Probability = 0
	Certain    Probability = 1
)

func (p Probability) Valid() bool {
	return p >= 0 && p <= 1
}

// --- }}}

// --- Discrete Distribution --- {{{

func NewDiscreteDistribution(d set.Interface) DiscreteDistribution {
	return &distribution{
		domain:   d,
		outcomes: set.New(),
		support:  make(map[Outcome]Probability),
	}
}

func NewUniformDiscrete(domain set.Interface) DiscreteDistribution {
	d := NewDiscreteDistribution(domain)

	individualSupport := Certain / Probability(domain.Cardinality())

	for o := range domain.Iter() {
		d.AddOutcome(o, individualSupport)
	}

	return d
}

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
	o := make(Outcomes, len(d.support))

	i := 0
	for k := range d.support {
		o[i] = k
	}

	return o
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

func Support(d Distribution) Probability {
	p := Probability(0)

	for o := range d.Outcomes().Iter() {
		p += d.ProbabilityOf(o)
	}

	return p
}

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

func Independent(d Distribution, X, Y RandomVariable) bool {
	return Covariance(d, X, Y) == 0
}
