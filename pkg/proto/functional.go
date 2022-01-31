package proto

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
)

// Shuffle implements Fisher-Yates shuffling.
func Shuffle(a []int) []int {
	var j int

	for i := len(a) - 1; i > 0; i-- {
		j = mathrand.Intn(i)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// SHuffleBig implements Fisher-Yates shuffling for big.Int.
func ShuffleBig(a []*big.Int) []*big.Int {
	var j int

	for i := len(a) - 1; i > 0; i-- {
		j = mathrand.Intn(i)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// NaiveIntersect returns the list of indices shared by two fixedlength
// arrays.
func NaiveIntersect(a []*big.Int, b []*big.Int) []int {

	intersection := make([]int, 0)

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].Cmp(b[j]) == 0 {
				intersection = append(intersection, i)
			}
		}
	}
	return intersection
}

// RandomOfBits returns a random number exactly bits long.
func RandomOfBits(maxRange *big.Int, bits int) (*big.Int, error) {
	r, err := rand.Int(rand.Reader, maxRange)
	for err != nil || r.BitLen() != 1024 {
		r, err = rand.Int(rand.Reader, maxRange)
	}
	return r, err
}
