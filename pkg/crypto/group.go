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
	ModExp(*big.Int, *big.Int) *big.Int
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

// ModExp computes modular exponent in the group.
// computes h(data)*k as in the hash is interpreted as a
// number and the exponent is calculated (mod exponentiation is
// used relative to the prime group used)
// TODO:Move the operation to the Group interface
func (pg *PrimeGroup) ModExpHashed(data *big.Int, exp *big.Int) *big.Int {
	modulus := pg.Order()
	h := HashBig(data)
	h_ := new(big.Int).SetBytes(h)
	h_ = new(big.Int).Mod(h_, modulus)
	modexp := new(big.Int).Exp(h_, exp, modulus)
	return modexp
}

// ModExp implements modular exponentiation in the group.
func (pg *PrimeGroup) ModExp(a *big.Int, k *big.Int) *big.Int {
	modulus := pg.Order()
	modexp := new(big.Int).Exp(a, k, modulus)
	return modexp
}

// GenPrimeGroup returns a prime order group with a random
// scalar generator.
func GenPrimeGroup(order *big.Int) Group {
	return &PrimeGroup{
		order:  order,
		random: rand.Reader,
	}
}
