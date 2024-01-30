package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func NewAesCrypt(key []byte) (Crypt, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AesCrypt{
		CipherBlock: cipherBlock,
		BlockSize:   aes.BlockSize,
	}, nil
}

type AesCrypt struct {
	CipherBlock cipher.Block
	BlockSize   int
}

func (this AesCrypt) Encrypt(content []byte) ([]byte, error) {
	contentBytes := pad([]byte(content), this.BlockSize)
	cipherText := make([]byte, this.BlockSize+len(contentBytes))
	iv := cipherText[:this.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(this.CipherBlock, iv)
	blockMode.CryptBlocks(cipherText[this.BlockSize:], contentBytes)
	return cipherText, nil
}

func (this AesCrypt) EncryptString(content string) (string, error) {
	encrypted, err := this.Encrypt([]byte(content))
	if err != nil {
		return "", err
	}

	b64String := base64.StdEncoding.EncodeToString(encrypted)
	return b64String, nil
}

func (this AesCrypt) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) < this.BlockSize || len(cipherText)%this.BlockSize != 0 {
		return nil, fmt.Errorf("invalid cipherText length")
	}

	iv := cipherText[:this.BlockSize]
	blockMode := cipher.NewCBCDecrypter(this.CipherBlock, iv)

	cipherText = cipherText[this.BlockSize:]
	blockMode.CryptBlocks(cipherText, cipherText)

	cipherText, err := unpad(cipherText)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func (this AesCrypt) DecryptString(b64String string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return "", err
	}

	decrypted, err := this.Decrypt(cipherText)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func pad(content []byte, blockSize int) []byte {
	pad := blockSize - len(content)%blockSize
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(content, padText...)
}

func unpad(content []byte) ([]byte, error) {
	contentSize := len(content)
	if contentSize == 0 {
		return nil, errors.New("content can't be empty")
	}

	pad := int(content[contentSize-1])
	if pad > contentSize || pad == 0 {
		return nil, errors.New("invalid padding")
	}

	//to verify, perhaps ok to be removed ?
	for _, v := range content[contentSize-pad:] {
		if int(v) != pad {
			return nil, errors.New("invalid padding character")
		}
	}

	return content[:contentSize-pad], nil
}
