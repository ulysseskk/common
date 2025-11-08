package passwordUtil

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowercase    = "abcdefghijklmnopqrstuvwxyz"
	uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}|;:',.<>/?"
)

func GeneratePassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	// Ensure the password has at least one character from each required set
	password := make([]byte, length)
	charsets := []string{lowercase, uppercase, digits, specialChars}

	// Generate one character from each charset and place them in random positions in the password
	for i, charset := range charsets {
		char, err := randomCharFromCharset(charset)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Fill the remaining password characters with random characters from all sets combined
	allChars := lowercase + uppercase + digits + specialChars
	for i := len(charsets); i < length; i++ {
		char, err := randomCharFromCharset(allChars)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Shuffle the password to ensure the required characters are not in predictable positions
	shuffledPassword, err := shufflePassword(password)
	if err != nil {
		return "", err
	}

	return string(shuffledPassword), nil
}

func randomCharFromCharset(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

func shufflePassword(password []byte) ([]byte, error) {
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(password))))
		if err != nil {
			return nil, err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}
	return password, nil
}
