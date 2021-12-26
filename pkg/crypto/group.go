package crypto

import (
	"crypto/rand"
	"io"
	"math/big"
)

// Group provides an abstraction to various groups.
type Group interface {
	Order() *big.Int
	RandomScalar() *big.Int
}

// PrimeGroup implements the group interface for
// prime order group over the integers.
type PrimeGroup struct {
	order  *big.Int
	random io.Reader
}

// Order returns the prime group order.
func (pg *PrimeGroup) Order() *big.Int {
	return pg.order
}

// RandomScalar returns a random scalar in the prime group.
// TODO:Can we replace calls to RandomScalar by a Fiat-Shamir transform ?
func (pg *PrimeGroup) RandomScalar() *big.Int {
	s, err := rand.Int(rand.Reader, pg.order)
	if err != nil {
		panic("failed to generate a proper random scalar")
	}
	return s
}

// GenPrimeGroup returns a prime order group with a random
// scalar generator.
func GenPrimeGroup(order *big.Int) Group {
	return &PrimeGroup{
		order:  order,
		random: rand.Reader,
	}
}
