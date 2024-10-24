package crypto

import (
	"github.com/thk-im/thk-im-base-server/utils"
	"strings"
)

type Crypto interface {
	Encrypt(text []byte) (string, error)
	Decrypt(cipherText string) ([]byte, error)
	EncryptUriBody(uri string, text []byte) (string, error)
	DecryptUriBody(uri string, cipherText string) ([]byte, error)
}

type defaultCrypto struct {
	aes       *utils.AES
	whiteList []string
}

func (d defaultCrypto) Encrypt(text []byte) (string, error) {
	return d.aes.Encrypt(text)
}

func (d defaultCrypto) Decrypt(cipherText string) ([]byte, error) {
	return d.aes.Decrypt(cipherText)
}

func (d defaultCrypto) EncryptUriBody(uri string, text []byte) (string, error) {
	for _, w := range d.whiteList {
		if len(w) > 0 && strings.HasPrefix(uri, w) {
			return string(text), nil
		}
	}
	return d.aes.Encrypt(text)
}

func (d defaultCrypto) DecryptUriBody(uri string, cipherText string) ([]byte, error) {
	for _, w := range d.whiteList {
		if len(w) > 0 && strings.HasPrefix(uri, w) {
			return []byte(cipherText), nil
		}
	}
	return d.aes.Decrypt(cipherText)
}

func NewCrypto(key, iv string, whiteList string) Crypto {
	tmp := strings.ReplaceAll(whiteList, " ", "")
	uriWhites := make([]string, 0)
	if len(tmp) > 0 {
		uriWhites = strings.Split(tmp, ",")
	}
	return &defaultCrypto{aes: utils.NewAES(key, iv), whiteList: uriWhites}
}
