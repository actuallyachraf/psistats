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
		bValues := []int{11, 13, 17, 21, 23, 29, 31, 37}

		// Setup
		// Group aggrement
		prime, _ := new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)
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

		// Sample r
		maxRange, _ := new(big.Int).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137216", 10)
		r, err := rand.Int(rand.Reader, maxRange)
		for err != nil || r.BitLen() != 1024 {
			r, err = rand.Int(rand.Reader, maxRange)
		}
		k := int64(len(I))

		// 0 <= r1 <= 2**128 - 1
		r1UpperBound, _ := new(big.Int).SetString(R1ScalarUpperBound, 10)

		// 2**511 <= r2 <= 2**512-1
		r2LowerBound, _ := new(big.Int).SetString(
			R2ScalarLowerBound, 10)
		r2UpperBound, _ := new(big.Int).SetString(R2ScalarUpperBound, 10)

		// Sample r1
		r1, err := rand.Int(rand.Reader, r1UpperBound)
		if err != nil {
			t.Fatal("failed to sample r1 with error :", err)
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
		t.Log("[+] Found valid r1 values : ", r1)

		r2, _ := rand.Int(rand.Reader, r2UpperBound)

		for r2.Cmp(r2LowerBound) < 0 {
			r2, _ = rand.Int(rand.Reader, r2UpperBound)
		}
		t.Log("[+] Found valid r2 values : ", r2)

		t.Log("[+] Additive Homomorphic Computation :")

		//sumTi := new(big.Int).SetInt64(0)
		// for _, idx := range I {
		// 	if idx+1 == len(encBvalues) {
		// 		t.Log("[+] Possible overflow returning...")
		// 		break
		// 	}
		// 	tmp := gaillier.Add(Bpub, encBvalues[idx].Bytes(), encBvalues[idx+1].Bytes())
		// 	sumTi.Add(sumTi, new(big.Int).SetBytes(tmp))
		// }
		//for _, idx := range I {
		//
		//	sumTi.Add(sumTi, encBvalues[idx])
		//}
		idx0 := I[0]
		idx1 := I[1]
		idx2 := I[2]
		c1 := encBvalues[idx0]
		c2 := encBvalues[idx1]
		c3 := encBvalues[idx2]

		c4 := gaillier.Add(Bpub, c1.Bytes(), c2.Bytes())
		c5 := gaillier.Add(Bpub, c4, c3.Bytes())

		constant := new(big.Int).Sub(r, r1)
		constant.Div(constant, bigK)

		rhs := gaillier.Mul(Bpub, c5, constant.Bytes())

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
		statistic := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).SetBytes(decrypted)), new(big.Float).SetInt(r))
		t.Log("[+] Computed Summary Statistic :", statistic)
		t.Log("[+] Plaintext Computed Summary Statistic :", (bValues[2]+bValues[4]+bValues[5])/len(I))

		statistic64, _ := statistic.Float64()
		expectedStatistic := (bValues[2] + bValues[4] + bValues[5]) / len(I)

		if math.Floor(statistic64) != float64(expectedStatistic) {
			t.Errorf("failed to assert output statistic consistency expected %d got %f", expectedStatistic, statistic64)
		}

	})
}
