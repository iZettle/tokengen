# tokengen

[![Build Status](https://travis-ci.org/intelligentpos/tokengen.svg?branch=master)](https://travis-ci.org/intelligentpos/tokengen)
[![Coverage Status](https://coveralls.io/repos/github/intelligentpos/tokengen/badge.svg?branch=refactor%2FECO-15)](https://coveralls.io/github/intelligentpos/tokengen?branch=refactor%2FECO-15)
[![GoDoc](https://godoc.org/github.com/intelligentpos/tokengen?status.svg)](https://godoc.org/github.com/intelligentpos/tokengen)
[![Go Report Card](https://goreportcard.com/badge/github.com/intelligentpos/tokengen)](https://goreportcard.com/report/github.com/intelligentpos/tokengen)

tokengen is small, simple and flexible token generator. tokengen allows you to specify
your character set and token length, and as such is ideally suited for generating secure
tokens in any language, random urls, passwords, and access tokens.

tokengen relies on the `crypto/rand` package, mapping values evenly to the character set
given, disregarding any values out of range.

```go

func GenerateOneTimePassword() (string, error){
    tokengen, err := tokengen.New(tokengen.DefaultCharset, 40)
    if err != nil {
        return tokengen, err
    }
    return tokengen.GenerateToken()
}

```

Please make sure that the character set and length of token you choose are large enough to ensure 
a reasonable amount of entropy.
