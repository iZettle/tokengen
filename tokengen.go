// Package tokengen provides a simple way to generate secure tokens of any
// length from any character set. Allowing for easy password, url, and token
// generation.
package tokengen

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
)

const (
	Base62         = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
	Base64         = Base62 + `+/`
	DefaultCharset = Base62
)

// TokenProvider is a simple interface for abstracting token provisioning.

type TokenProvider interface {
	GenerateToken() (string, error)
}

// Tokengen implements TokenProvider and contains the configuration for
// generating cryptographically secure tokens.
type Tokengen struct {
	Charset string
	Length  int
}

// GenerateToken will provide a string of letters, picked at random from
// the given character set, with even distribution of runes from the set
// GenerateToken relies on the crypto/rand package for it's random data
// source, rather than the math package, so is ideally suited for secure
// uses such as password, token and url generation.
//  func GenerateOneTimePassword() (string, error){
//  	tokengen := tokengen.Tokengen{
//  		Length: 40,
//  		Charset: tokengen.DefaultCharset,
//  	}
//  	return tokengen.GenerateToken()
//  }
func (t Tokengen) GenerateToken() (string, error) {
	if t.Length < 0 {
		return "", errors.New("length cannot be negative")
	}
	if len(t.Charset) == 0 {
		return "", errors.New("character set too small")
	}
	runes, err := t.generateToken(t.Length)
	return string(runes), err
}

func (t Tokengen) generateToken(length int) ([]rune, error) {
	charset := []rune(t.Charset)

	// find the minimum number of bytes the index can be represented in,
	bytesPerRune, maxIndex := bytesPerRuneIndex(len(charset))

	// the max value within the given range that len(charset) % x == 0
	maxValue := maxIndex - (maxIndex % len(charset))

	// extend set by amount percentage chance that value is out of range
	requiredBytes := t.calcRequiredBytes(maxIndex, maxValue, bytesPerRune)
	randomBytes := make([]byte, requiredBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	output := make([]rune, 0, length)
	padding := make([]byte, 4-bytesPerRune)

	for i := 0; i < requiredBytes && len(output) != length; {
		// take a slice of bytes of the size of the index
		idxBytes := randomBytes[i : i+bytesPerRune]

		//pad them to fit a uint32
		idxBytes = append(padding, idxBytes...)
		randValue := binary.BigEndian.Uint32(idxBytes)

		// if out of range, move ahead one byte
		if randValue >= uint32(maxValue) {
			i++
			continue
		}
		idx := int(randValue) % len(charset)
		output = append(output, charset[idx])
		i += bytesPerRune
	}
	if len(output) == length {
		return output, nil
	}

	// we may not have had enough random data in range, call again.
	runes, err := t.generateToken(length - len(output))
	if err != nil {
		return nil, err
	}
	output = append(output, runes...)
	return output, nil
}

func bytesPerRuneIndex(numRunes int) (bytesPerRune, maxValue int) {
	maxValue = 1
	// iterate until we find a multiple of 256 greater than the number of runes
	for bytesPerRune = 0; maxValue <= numRunes; bytesPerRune++ {
		maxValue *= (1 << 8)
	}
	maxValue -= 1
	return
}

func (t *Tokengen) calcRequiredBytes(maxIndex, maxValue, bytesPerRune int) int {
	multiplier := float64(maxIndex) / float64(maxValue)
	increasedLength := math.Ceil(multiplier * float64(t.Length))
	return int(increasedLength) * bytesPerRune
}
