package proto

import (
	"crypto/rand"
	"math/big"

	"github.com/actuallyachraf/psistats/pkg/crypto"
)

// HashIDs computes the SHA-256 of a party's IDs.
func HashIDs(ids []int, k *big.Int, G crypto.Group) []*big.Int {

	hashedIds := make([]*big.Int, len(ids))
	for i, id := range ids {
		valBig := new(big.Int).SetInt64(int64(id))
		// H(id)
		hashedId := crypto.HashToBig(valBig)
		// H(id)^k
		modexp := G.ModExp(hashedId, k)
		hashedIds[i] = modexp
	}
	return hashedIds
}

// ExpHashIDs lifts hashed IDs to an exponent this is done by both
// parties.
func ExpHashIDs(hashedIds []*big.Int, k *big.Int, G crypto.Group) []*big.Int {
	liftedHashedIds := make([]*big.Int, len(hashedIds))

	for i, id := range hashedIds {
		modexp := G.ModExp(id, k)
		liftedHashedIds[i] = modexp
	}
	return liftedHashedIds
}

// SampleR1 will sample R1 congruent to r mod k where k is the
// cardinality of the intersection.
func SampleR1(r *big.Int, k int64) (*big.Int, error) {
	// Sample r1
	r1UpperBound, _ := new(big.Int).SetString(R1ScalarUpperBound, 10)

	r1, err := rand.Int(rand.Reader, r1UpperBound)
	if err != nil {
		return nil, err
	}
	res := new(big.Int).SetInt64(0)
	diff := new(big.Int).Sub(r1, r)
	zero := new(big.Int).SetInt64(0)
	bigK := new(big.Int).SetInt64(k)

	res.Mod(diff, bigK)
	for res.Cmp(zero) != 0 {
		r1, _ = rand.Int(rand.Reader, r1UpperBound)
		diff.Sub(r1, r)
		res.Mod(diff, bigK)
	}
	return r1, nil
}

// SampleR2 will sample R2 with r2 in range [2**511,2**512 - 1].
func SampleR2() *big.Int {
	r2UpperBound, _ := new(big.Int).SetString(R2ScalarUpperBound, 10)
	r2LowerBound, _ := new(big.Int).SetString(R2ScalarLowerBound, 10)
	r2, _ := rand.Int(rand.Reader, r2UpperBound)

	for r2.Cmp(r2LowerBound) < 0 {
		r2, _ = rand.Int(rand.Reader, r2UpperBound)
	}
	return r2
}

// SampleR will sample a random number R of exactly 1024-bits long.
func SampleR() (*big.Int, error) {
	maxRange, _ := new(big.Int).SetString(RScalarUpperBound, 10)
	r, err := RandomOfBits(maxRange, 1024)
	if err != nil {
		return nil, err
	}
	return r, nil
}
