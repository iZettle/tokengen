// Package tokengen provides a simple way to generate secure tokens of any
// length from any character set. Allowing for easy password, url, and token
// generation.
package tokengen

import (
	"crypto/rand"
	"errors"
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

func New(charset string, length int) (Tokengen, error) {
	if length < 1 {
		return Tokengen{}, errors.New(`length must be a positive integer`)
	}
	if len(charset) == 0 {
		return Tokengen{}, errors.New(`charset must contain characters`)
	}
	return Tokengen{
		Charset:     charset,
		Length:      length,
		distributor: newRuneDistributor([]rune(charset), length, rand.Reader),
	}, nil
}

// Tokengen implements TokenProvider and contains the configuration for
// generating cryptographically secure tokens.
type Tokengen struct {
	Charset     string
	Length      int
	distributor runeDistributor
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
	runes, err := t.distributor.generateToken()
	if err != nil {
		return ``, err
	}
	return string(runes), nil
}
