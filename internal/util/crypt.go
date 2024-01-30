package util

type Crypt interface {
	Encrypt(content []byte) ([]byte, error)
	EncryptString(content string) (string, error)

	Decrypt(cipherText []byte) ([]byte, error)
	DecryptString(b64string string) (string, error)
}
