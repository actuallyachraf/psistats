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
	})
}
