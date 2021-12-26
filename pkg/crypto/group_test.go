package crypto

import (
	"math/big"
	"testing"

	"github.com/actuallyachraf/psistats/pkg/proto"
)

func TestGroup(t *testing.T) {
	t.Run("TestPrimeGroup", func(t *testing.T) {
		order, ok := new(big.Int).SetString(proto.MODPGroup1536, 16)
		if !ok {
			t.Fatal("TestPrimeGroup failed with error : bad string")
		}
		pg := GenPrimeGroup(order)

		isLessThan := func(a *big.Int, b *big.Int) bool {
			return a.Cmp(b) == -1
		}

		for i := 0; i < 10; i++ {
			t.Log(pg.RandomScalar())
			if !isLessThan(pg.RandomScalar(), order) {
				t.Fatal("TestPrimeGroup failed random scalar is higher than order")
			}
		}
	})
}
