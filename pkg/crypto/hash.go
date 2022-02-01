package crypto

import (
	"crypto/sha256"
	"math/big"
)

// H defines a hash function in this case SHA-256.
func H(data []byte) []byte {
	b := sha256.Sum256(data)
	return b[:]
}

// HashBig returns the SHA-256 of a big integer.
func HashBig(x *big.Int) []byte {
	return H(x.Bytes())
}

// HashToBig returns the SHA-256 but casted to a big.Int
func HashToBig(x *big.Int) *big.Int {
	return new(big.Int).SetBytes(HashBig(x))
}
