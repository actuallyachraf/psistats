package crypto

import (
	"math/big"
	"testing"
)

const MODPGroup1536 = `FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA237327FFFFFFFFFFFFFFFF`

func TestGroup(t *testing.T) {
	t.Run("TestPrimeGroup", func(t *testing.T) {

		assertOrderIsPrime := func(order *big.Int) bool {
			return order.ProbablyPrime(256)
		}

		order, ok := new(big.Int).SetString(MODPGroup1536, 16)
		if !ok || !assertOrderIsPrime(order) {
			t.Fatal("TestPrimeGroup failed with error : order is not prime")
		}
		pg := GenPrimeGroup(order)
		if !assertOrderIsPrime(pg.Order()) {
			t.Fatal("TestPrimeGroup failed with error : order is not prime")
		}
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
