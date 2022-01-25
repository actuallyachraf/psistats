package crypto

import (
	"crypto/sha256"
	"encoding/binary"
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

// HashInt64 returns the SHA-256 of a <= 64-bit signed integer.
// Integers are encoded using Go's proper Variable length encoding.
// TODO:might be more sane to replace VarInt with a standard encoding
// maybe LEB128 ?
func HashInt64(x int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, x)
	return H(buf)
}

// HashToBig returns the SHA-256 but casted to a big.Int
func HashToBig(x *big.Int) *big.Int {
	return new(big.Int).SetBytes(HashBig(x))
}
