package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

func getGCMSalt(gcm cipher.AEAD) []byte {
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil
	}
	return nonce
}

func EncryptGCM(plaintext string) ([]byte, error) {
	secretKey := os.Getenv("DB_SECRET_KEY")

	if secretKey == "" {
		panic("Unabled to get the db secret key from .env")
	}

	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	text := []byte(plaintext)

	salt := getGCMSalt(gcm)
	if salt == nil {
		return nil, errors.New("error while to generate a new salt")
	}

	ciphertext := gcm.Seal(salt, salt, text, nil)

	return ciphertext, nil
}

func DecryptGCM(ciphertext string) (string, error) {
	secretKey := os.Getenv("DB_SECRET_KEY")

	if secretKey == "" {
		panic("Unabled to get the db secret key from .env")
	}

	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func HashArgon2(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	encodedSalt := fmt.Sprintf("%x", salt)
	encodedHashedPassword := fmt.Sprintf("%x", hashedPassword)

	return encodedSalt + ":" + encodedHashedPassword, nil
}

func VerifyArgon2(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, ":")
	if len(parts) != 2 {
		return false
	}

	salt, _ := hex.DecodeString(parts[0])
	hashedPassword, _ := hex.DecodeString(parts[1])

	calculatedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return bytes.Equal(hashedPassword, calculatedHash)
}
