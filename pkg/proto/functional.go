package proto

import (
	"math/big"
	"math/rand"
)

// Shuffle implements Fisher-Yates shuffling.
func Shuffle(a []int) []int {
	var j int

	for i := len(a) - 1; i > 0; i-- {
		j = rand.Intn(i)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// SHuffleBig implements Fisher-Yates shuffling for big.Int.
func ShuffleBig(a []*big.Int) []*big.Int {
	var j int

	for i := len(a) - 1; i > 0; i-- {
		j = rand.Intn(i)
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
				intersection = append(intersection, j)
			}
		}
	}
	return intersection
}
