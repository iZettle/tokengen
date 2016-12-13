# tokengen

tokengen is small, simple and flexible token generator. tokengen allows you to specify
your character set and token length, and as such is ideally suited for generating secure
tokens in any language, random urls, passwords, and access tokens.

tokengen relies on the `crypto/rand` package, mapping values evenly to the character set
given, disregarding any values out of range.

```go

func GenerateOneTimePassword() (string, error){
    tokengen := tokengen.Tokengen{
        Length: 40,
        Charset: tokengen.DefaultCharset,
    }
    return tokengen.GenerateToken()
}

```

Please make sure that the character set and length of token you choose are large enough to ensure 
a reasonable amount of entropy.
