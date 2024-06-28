package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AES struct {
	key []byte
	iv  []byte
}

func NewAES(key string, iv string) *AES {
	return &AES{
		key: []byte(key),
		iv:  []byte(iv),
	}
}

func (a *AES) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	paddedPlaintext := a.pkcs7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, a.iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *AES) Decrypt(ciphertext string) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	decryptedData := make([]byte, len(decodedCiphertext))
	mode := cipher.NewCBCDecrypter(block, a.iv)
	mode.CryptBlocks(decryptedData, decodedCiphertext)
	return a.pkcs7UnPadding(decryptedData), nil
}

// 使用PKCS7填充方式对数据进行填充
func (a *AES) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 对使用PKCS7填充方式的数据进行去填充
func (a *AES) pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	if length < 1 {
		return []byte{}
	}
	unPadding := int(data[length-1])
	if length < unPadding {
		return []byte{}
	}
	return data[:(length - unPadding)]
}
