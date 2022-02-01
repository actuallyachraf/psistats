package crypto

import (
	"errors"
	"math/big"

	"github.com/actuallyachraf/gomorph/gaillier"
)

// Encrypt a message.
func Encrypt(pub *gaillier.PubKey, message []byte) ([]byte, error) {
	return gaillier.Encrypt(pub, message)

}

// EncryptInt64 encryptes multiple int64 messages.
func EncryptInt64(pub *gaillier.PubKey, messages []int64) ([]*big.Int, error) {
	ciphers := make([]*big.Int, len(messages))
	for i, message := range messages {
		messageBig := new(big.Int).SetInt64(message)
		cipher, err := Encrypt(pub, messageBig.Bytes())
		if err != nil {
			return nil, err
		}
		ciphers[i] = new(big.Int).SetBytes(cipher)
	}
	return ciphers, nil
}

// Dec decrypts a Paillier encrypted message.
func Decrypt(priv *gaillier.PrivKey, enc []byte) ([]byte, error) {
	return gaillier.Decrypt(priv, enc)
}

// AddEnc adds up encrypted values.
func AddEnc(pub *gaillier.PubKey, values ...[]byte) ([]byte, error) {
	if len(values) < 2 {
		return nil, errors.New("expected values to be atleast two")
	}

	cipherOut := values[0]
	for i := 1; i < len(values); i++ {
		cipherOut = gaillier.Add(pub, cipherOut, values[i])
	}

	return cipherOut, nil
}
