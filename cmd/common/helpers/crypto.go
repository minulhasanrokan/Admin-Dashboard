package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func getEncryptionKey() []byte {

	key := []byte(os.Getenv("APP_KEY"))

	if len(key) == 0 {

		key = []byte("default-secret-key-change-this-32b")
	}

	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		return padded
	}

	return key[:32]
}

func Encrypt(plainText string) (string, error) {

	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

func Decrypt(cipherHex string) (string, error) {

	if cipherHex == "" {
		return "", nil
	}

	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(cipherText) < nonceSize {
		return "", errors.New("invalid ciphertext format")
	}

	nonce, cipherTextBytes := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func EncryptPtr(val string) *string {
	if val == "" {
		return nil
	}
	encrypted, err := Encrypt(val)
	if err != nil {
		return nil
	}
	return &encrypted
}

func DecryptPtr(val *string) *string {

	if val == nil || *val == "" {
		return nil
	}

	decrypted, err := Decrypt(*val)

	if err != nil {
		return val
	}

	return &decrypted
}
