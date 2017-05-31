package tokengen

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type dummyReader struct{}

func (dr dummyReader) Read(_ []byte) (int, error) {
	return 0, errors.New(`could not read`)
}

func TestTokengen_GenerateTokenMultipleLengths(t *testing.T) {

	testCases := []int{
		32,
		4096,
		64,
		12,
		1,
	}
	for _, expected := range testCases {
		t.Run(fmt.Sprintf(`Lenght: %v`, expected), func(t *testing.T) {
			tg, err := New(DefaultCharset, expected)
			if err != nil {
				t.Fatal(err)
			}
			token, err := tg.GenerateToken()
			if actual := len(token); actual != expected {
				t.Fatalf("Expected length of %v got %v for tokengen of %v", expected, actual, tg)
			}
		})
	}
}

func TestTokengen_OccurenceWithinTwoPercentOfAverage(t *testing.T) {
	if testing.Short() {
		t.Skip(`this test should be ran if you've done a major refactor.'`)
	}
	file, err := os.Open(`random.data`)
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	tokengen := Tokengen{
		distributor: newRuneDistributor([]rune(DefaultCharset), 1<<12, file),
	}
	scores := genEmptyMap(DefaultCharset)
	for i := 0; i < (1 << 9); i++ {
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

	avg := total / len(DefaultCharset)
	twoPercent := avg / 50
	distribution := map[string]int{}
	for char, val := range scores {
		dev := val - avg
		if dev < 0 {
			dev = -dev
		}
		distribution[char] = dev
		if dev > twoPercent {
			t.Fatalf(`there is a greater than two percent deviaton from average of the character '%s' with a deviation of %v compared to a two percent being %v`, char, dev, twoPercent)
		}
	}
}

func TestTokengen_InvalidLength(t *testing.T) {
	_, err := New(DefaultCharset, -13)
	if err == nil {
		t.Fatal(`no error returned for error case`)
	}
}

func TestTokengen_InvalidCharset(t *testing.T) {
	_, err := New(``, 40)
	if err == nil {
		t.Fatal(`no error returned for error case`)
	}
}

func TestTokengen_GenerateToken(t *testing.T) {
	tokengen, err := New(`0123456789abcdef`, 64)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tokengen.GenerateToken(); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestTokengen_GenerateTokenErroringRandSource(t *testing.T) {
	tokengen := Tokengen{
		distributor: newRuneDistributor([]rune(DefaultCharset), 1<<12, dummyReader{}),
	}
	_, err := tokengen.GenerateToken()
	if err == nil {
		t.Fatal(`no error for invalid randsource`)
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
