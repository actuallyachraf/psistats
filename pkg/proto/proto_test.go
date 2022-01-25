package proto

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/actuallyachraf/gomorph/gaillier"
	"github.com/actuallyachraf/psistats/pkg/crypto"
)

func TestProtocolPrimitives(t *testing.T) {
	t.Run("TestShuffle", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		aShuffled := Shuffle(a)
		t.Log(aShuffled)
	})
}

func TestProtocolInstance(t *testing.T) {
	t.Run("TestPSIMean-SingleInstance", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5, 6}
		b := []int{0, 0, 3, 0, 5, 6}
		bValues := []int{11, 13, 17, 21, 23, 29}

		// Setup
		// Group aggrement
		prime, _ := new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)
		G := crypto.GenPrimeGroup(prime)
		// A keypair
		_, _, err := gaillier.GenerateKeyPair(rand.Reader, 128)
		if err != nil {
			t.Fatal("[Alice]failed to generate keypair with error :", err)
		}
		// B keypair
		Bpub, _, err := gaillier.GenerateKeyPair(rand.Reader, 128)
		if err != nil {
			t.Fatal("[Bob]failed to generate keypair with error :", err)
		}

		// A's encryption phase
		encA := make([]*big.Int, 0)
		k1 := G.RandomScalar()
		for _, val := range a {
			// Compute h(ai)^k1
			valBig := new(big.Int).SetInt64(int64(val))
			// H(ai)
			hashed := crypto.HashToBig(valBig)
			// H(ai)^k1 mod G
			modexp := G.ModExp(hashed, k1)
			encA = append(encA, modexp)
		}

		// B's encryption phase
		encBshared := make([]*big.Int, 0)
		k2 := G.RandomScalar()
		// Compute h(ai)^k1k2 <=> (h(ai)^k1)^k2
		for _, val := range encA {
			modexp := G.ModExp(val, k2)
			encBshared = append(encBshared, modexp)
		}
		// Compute {h(bj)^k2, E(tj)}
		// Essentially Compute h(bj)^k2
		// Encrypt B's values using his private key
		encB := make([]*big.Int, 0)
		encBvalues := make([]*big.Int, 0)
		for _, val := range b {
			valBig := new(big.Int).SetInt64(int64(val))
			hashed := crypto.HashToBig(valBig)
			modexp := G.ModExp(hashed, k2)
			encB = append(encB, modexp)
		}
		for _, val := range bValues {
			valBig := new(big.Int).SetInt64(int64(val))
			enc, err := gaillier.Encrypt(Bpub, valBig.Bytes())
			if err != nil {
				t.Fatal("[Bob]failed to encrypt B values using the Gaillier instance with error: ", err)
			}
			encBvalues = append(encBvalues, new(big.Int).SetBytes(enc))
		}
		t.Log("[Bob] E(tj) : ", encBvalues)
		// B shuffles h(ai)k1k2
		shuffledBshared := ShuffleBig(encBshared)
		// B sends {h(bj)^k2,E(tj)} to A
		// encB is h(bj)^k2
		// encBvalues is E(tj)
		// encBshared is h(ai)^k1k2

		// A's match and compute
		// A computes h(bj)^k1k2 by lifting (h(bj)k2)^k1
		encAshared := make([]*big.Int, 0)

		for _, val := range encB {
			modexp := G.ModExp(val, k1)
			encAshared = append(encAshared, modexp)
		}
		//
		// Compute the set intersection I as:
		// I = {j : h(ai)^k1k2 == h(bj)^k2k1}
		// I is the set of B's indices that are in the interesection
		t.Log("[Alice]h(bj)^k2k1 : ", encAshared)
		t.Log("[Bob]h(ai)^k1k2 :", shuffledBshared)
		I := NaiveIntersect(encAshared, shuffledBshared)

		t.Log("Alice computed PSI : ", I)
	})
}
