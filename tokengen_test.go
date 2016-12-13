package tokengen

import (
	"math"
	"testing"
)

func TestTokengen_GenerateToken(t *testing.T) {
	cases := []int{32, 4096, 64, 12, 1}

	tg := Tokengen{Charset: DefaultCharset}
	for _, expected := range cases {
		tg.Length = expected
		if token, err := tg.GenerateToken(); err == nil {
			if actual := len(token); actual != expected {
				t.Logf("Expected length of %v got %v for tokengen of %v", expected, actual, tg)
				t.FailNow()
			}
		} else {
			t.Log(err)
			t.FailNow()
		}
	}
}

func TestTokengen_GenerateToken2(t *testing.T) {
	tokengen := Tokengen{
		Charset: DefaultCharset,
		Length:  (1 << 12),
	}
	scores := genEmptyMap(DefaultCharset)
	for i := 0; i < (1 << 12); i++ {
		token, err := tokengen.GenerateToken()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		scores = total(token, scores)
	}
	total := 0
	for _, score := range scores {
		total += score
	}
	avg := float64(total) / float64(len(DefaultCharset))
	differences := map[string]float64{}

	min, max := 0.0, 0.0

	var minstr, maxstr string
	for key, score := range scores {
		res := 100 - ((avg / float64(score)) * 100)
		if res > max {
			maxstr = key
			max = res
		}
		if res < min {
			minstr = key
			min = res
		}
		differences[key] = res
	}
	if math.Abs(min)*1.15 < math.Abs(max) {
		t.Log(minstr, min, maxstr, max)
		t.FailNow()
	}
}

func TestTokengen_GenerateToken3(t *testing.T) {
	tokengen := Tokengen{
		Charset: DefaultCharset,
		Length:  -13,
	}

	if token, err := tokengen.GenerateToken(); err == nil {
		t.Log(tokengen, token)
		t.FailNow()
	}
}

func TestTokengen_GenerateToken4(t *testing.T) {
	tokengen := Tokengen{
		Charset: `0123456789abcdef`,
		Length:  64,
	}

	if _, err := tokengen.GenerateToken(); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func genEmptyMap(charset string) map[string]int {
	output := map[string]int{}
	for _, char := range charset {
		output[string(char)] = 0
	}
	return output
}

func total(token string, scores map[string]int) map[string]int {
	for _, char := range token {
		scores[string(char)]++
	}
	return scores
}
