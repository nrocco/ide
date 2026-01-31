package ide

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	upperChars       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars       = "abcdefghijklmnopqrstuvwxyz"
	digitChars       = "0123456789"
	punctuationChars = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	hexChars         = "0123456789abcdef"
)

// PasswordOptions can be used to control character types that are part of the password
type PasswordOptions struct {
	Upper       bool
	Lower       bool
	Numbers     bool
	Punctuation bool
	Hexdigits   bool
}

// Charset returns all the characters from which to generate a password
func (o PasswordOptions) Charset() string {
	if o.Hexdigits {
		return hexChars
	}

	var chars string
	if o.Upper {
		chars += upperChars
	}
	if o.Lower {
		chars += lowerChars
	}
	if o.Numbers {
		chars += digitChars
	}
	if o.Punctuation {
		chars += punctuationChars
	}
	return chars
}

// GeneratePassword generates a password with the given length and options
func GeneratePassword(length int, opts PasswordOptions) (string, error) {
	chars := opts.Charset()
	if len(chars) == 0 {
		return "", errors.New("no character set selected")
	}

	result := make([]byte, length)
	charsLen := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		result[i] = chars[idx.Int64()]
	}

	return string(result), nil
}
