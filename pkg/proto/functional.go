package proto

import "math/rand"

// functional.go provides a few useful algorithm.

// Shuffle implements modern Fisher-Yates shuffle.
func Shuffle(a []int) []int {
	var j int

	for i := len(a) - 1; i > 0; i-- {
		j = rand.Intn(i)
		a[i], a[j] = a[j], a[i]
	}
	return a
}
