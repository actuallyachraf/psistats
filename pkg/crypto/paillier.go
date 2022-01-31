package crypto

import (
	"github.com/actuallyachraf/gomorph/gaillier"
)

// Import Paillier
// Define a Class Paillier and instantiate the class
// Enc()
// Dec()
// Add()
// MulConst()
func Encrypt(pub *gaillier.PubKey, message []byte) ([]byte, error) {
	return gaillier.Encrypt(pub, message)

}

// Dec decrypts a Paillier encrypted message.
func Decrypt(priv *gaillier.PrivKey, enc []byte) ([]byte, error) {
	return gaillier.Decrypt(priv, enc)
}

// AddEnc adds up encrypted values.
func AddEnc(pub *gaillier.PubKey, values ...[]byte) []byte {
	// if len(values) != 2 {
	// 	return nil, errors.New("expected values to be of length two")
	// }

	cipherOut := values[0]
	for i := 1; i < len(values); i++ {
		cipherOut = gaillier.Add(pub, cipherOut, values[i])
	}

	return cipherOut
}
