package proto

import (
	"crypto/rand"
	"math"
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
		a := []int{1, 2, 3, 4, 5, 6, 8, 9}
		b := []int{0, 0, 3, 0, 5, 6, 9, 10}
		bValues := []int64{11, 13, 17, 21, 23, 29, 31, 37}

		// Setup
		// Group aggrement
		prime, _ := new(big.Int).SetString(MODPGroup256, 10)
		G := crypto.GenPrimeGroup(prime)
		// A keypair
		_, _, err := gaillier.GenerateKeyPair(rand.Reader, 2048)
		if err != nil {
			t.Fatal("[Alice]failed to generate keypair with error :", err)
		}
		// B keypair
		Bpub, Bpriv, err := gaillier.GenerateKeyPair(rand.Reader, 2048)
		if err != nil {
			t.Fatal("[Bob]failed to generate keypair with error :", err)
		}

		// A's encryption phase
		k1 := G.RandomScalar()
		encA := HashIDs(a, k1, G)
		// B's encryption phase
		k2 := G.RandomScalar()
		encBshared := ExpHashIDs(encA, k2, G)
		// Compute {h(bj)^k2, E(tj)}
		// Essentially Compute h(bj)^k2
		// Encrypt B's values using his private key
		encB := HashIDs(b, k2, G)
		encBvalues, err := crypto.EncryptInt64(Bpub, bValues)
		if err != nil {
			t.Fatal("[Bob]failed to encrypt B values using the Gaillier instance with error: ", err)
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
		encAshared := ExpHashIDs(encB, k1, G)
		//
		// Compute the set intersection I as:
		// I = {j : h(ai)^k1k2 == h(bj)^k2k1}
		// I is the set of B's indices that are in the interesection
		t.Log("[Alice]h(bj)^k2k1 : ", encAshared)
		t.Log("[Bob]h(ai)^k1k2 :", shuffledBshared)
		I := NaiveIntersect(encAshared, shuffledBshared)

		t.Log("Alice computed PSI : ", I)

		// Sample r
		r, err := SampleR()
		if err != nil {
			t.Fatal("[Alice] failed to sample a random number with error :", err)
		}
		k := int64(len(I))
		bigK := new(big.Int).SetInt64(k)

		// Sample r1
		r1, err := SampleR1(r, k)
		if err != nil {
			t.Fatal("[Alice] failed to sample correct R1 with error :", err)
		}
		t.Log("[+] Found valid r1 values : ", r1)

		r2 := SampleR2()
		t.Log("[+] Found valid r2 values : ", r2)

		t.Log("[+] Additive Homomorphic Computation :")

		encryptedValues := make([][]byte, 0)
		for _, idx := range I {
			encryptedValues = append(encryptedValues, encBvalues[idx].Bytes())
		}
		sum, err := crypto.AddEnc(Bpub, encryptedValues...)
		if err != nil {
			t.Fatal("[Alice] failed to sum values with error :", err)
		}
		constant := new(big.Int).Sub(r, r1)
		constant.Div(constant, bigK)

		rhs := gaillier.Mul(Bpub, sum, constant.Bytes())

		encr := gaillier.AddConstant(Bpub, rhs, r2.Bytes())

		t.Log("[+] Encrypted Message Size in Bits : ", len(encr)*8)

		// B's decryption phase
		decrypted, err := gaillier.Decrypt(Bpriv, encr)
		if err != nil {
			t.Fatal("[Bob] failed to decrypt encrypted value with error :", err)
		}
		t.Log("[+] R bitlen :", r.BitLen())
		t.Log("[+] Decrypted Message :", new(big.Int).SetBytes(decrypted))
		// Approximation
		mean := func(idx []int, values []int64) int {
			sum := int64(0)
			for _, i := range idx {
				sum += values[i]
			}
			return int(sum / int64(len(idx)))
		}
		statistic := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).SetBytes(decrypted)), new(big.Float).SetInt(r))
		t.Log("[+] Computed Summary Statistic :", statistic)
		t.Log("[+] Plaintext Computed Summary Statistic :", mean(I, bValues))

		statistic64, _ := statistic.Float64()
		expectedStatistic := mean(I, bValues)

		if math.Floor(statistic64) != float64(expectedStatistic) {
			t.Errorf("failed to assert output statistic consistency expected %d got %f", expectedStatistic, statistic64)
		}

	})

}
