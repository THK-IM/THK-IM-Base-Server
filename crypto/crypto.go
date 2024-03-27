package crypto

import "github.com/thk-im/thk-im-base-server/utils"

type Crypto interface {
	Encrypt(text []byte) (string, error)
	Decrypt(cipherText string) ([]byte, error)
}

type defaultCrypto struct {
	aes *utils.AES
}

func (d defaultCrypto) Encrypt(text []byte) (string, error) {
	return d.aes.Encrypt(text)
}

func (d defaultCrypto) Decrypt(cipherText string) ([]byte, error) {
	return d.aes.Decrypt(cipherText)
}

func NewCrypto(key, iv string) Crypto {
	return &defaultCrypto{aes: utils.NewAES(key, iv)}
}
