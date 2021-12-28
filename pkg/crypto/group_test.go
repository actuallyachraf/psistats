package crypto

import (
	"math/big"
	"testing"

	"github.com/actuallyachraf/psistats/pkg/proto"
)

func TestGroup(t *testing.T) {
	t.Run("TestPrimeGroup", func(t *testing.T) {

		assertOrderIsPrime := func(order *big.Int) bool {
			return order.ProbablyPrime(256)
		}

		order, ok := new(big.Int).SetString(proto.MODPGroup1536, 16)
		if !ok || !assertOrderIsPrime(order) {
			t.Fatal("TestPrimeGroup failed with error : order is not prime")
		}
		pg := GenPrimeGroup(order)

		isLessThan := func(a *big.Int, b *big.Int) bool {
			return a.Cmp(b) == -1
		}

		for i := 0; i < 10; i++ {
			if !isLessThan(pg.RandomScalar(), order) {
				t.Fatal("TestPrimeGroup failed random scalar is higher than order")
			}
		}
		// Intersection Mean
		// TODO: Separate the protocol tests into proper proto package
		// Step 1
		// select a random scalar k1 in G
		// compute h(a)**k
		k1 := pg.RandomScalar()
		x := new(big.Int).SetInt64(124314)
		secAlice := pg.ModExp(x, k1)
		// Step 2
		// select random k2 in G
		// compute h(a)**k1k2
		k2 := pg.RandomScalar()
		y := new(big.Int).SetInt64(123456)
		secBob := pg.ModExp(y, new(big.Int).Mul(k1, k2))
		// Step 3
		// Bob computes h(b)**k2
		// Bob computes Enc(t) t: associated statistic with identifier
		t.Log("Alice's secret : " + secAlice.Text(16))
		t.Log("Bob's secret : " + secBob.Text(16))

	})
}
