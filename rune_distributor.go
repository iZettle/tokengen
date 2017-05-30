package tokengen

import (
	"encoding/binary"
	"io"
	"math"
)

const intBytes = 1 << 2

type runeDistributor struct {
	runeset                         []rune
	randSource                      io.Reader
	bytsPerRuneIdx, maxValueOfBytes int
	setLength, throwawayLimit       int
	requiredBytes                   int
	padding                         []byte
}

func newRuneDistributor(charset []rune, length int, randSource io.Reader) runeDistributor {
	bytes, maxValue := bytesPerRuneIndex(len(charset))
	throwawayLimit := calcThrowawayLimit(maxValue, len(charset))
	return runeDistributor{
		runeset:         charset,
		randSource:      randSource,
		bytsPerRuneIdx:  bytes,
		maxValueOfBytes: maxValue,
		setLength:       length,
		throwawayLimit:  throwawayLimit,
		requiredBytes:   calcBytesRequired(length, maxValue, throwawayLimit, bytes),
		padding:         make([]byte, intBytes-bytes),
	}
}

func (rd runeDistributor) getRandomData() ([]byte, error) {
	randomBytes := make([]byte, rd.requiredBytes)
	if _, err := rd.randSource.Read(randomBytes); err != nil {
		return randomBytes, err
	}
	return randomBytes, nil
}

func (rd runeDistributor) generateToken() ([]rune, error) {
	output := make([]rune, 0, rd.setLength)
	randBytes, err := rd.getRandomData()
	if err != nil {
		return output, err
	}
	for i := 0; i < rd.requiredBytes && len(output) != rd.setLength; {
		idxBytes := append(rd.padding, randBytes[i:i+rd.bytsPerRuneIdx]...)
		randValue := binary.BigEndian.Uint32(idxBytes)
		if randValue >= uint32(rd.throwawayLimit) {
			i++
			continue
		}
		idx := int(randValue) % len(rd.runeset)
		output = append(output, rd.runeset[idx])
		i += rd.bytsPerRuneIdx
	}
	if len(output) == rd.setLength {
		return output, nil
	}
	extra, err := newRuneDistributor(rd.runeset, rd.setLength-len(output), rd.randSource).generateToken()
	if err != nil {
		return output, err
	}
	return append(output, extra...), nil
}

func bytesPerRuneIndex(numRunes int) (bytes, maxValue int) {
	const maxByteValue = 1 << 8
	permutations := 1

	// iterate until we find the number of bytes that
	// supports the number of runes we have (or more)
	for bytes = 0; permutations <= numRunes; bytes++ {
		permutations *= maxByteValue
	}

	// while the number of permutations of values is a for the bytes is this,
	// the max is one less, due to 0 being a value too.
	maxValue = permutations - 1
	return
}

func calcThrowawayLimit(maxIdx, charsetLen int) int {
	return maxIdx - (maxIdx % charsetLen)
}

func calcBytesRequired(length, maxIdx, throwawayLimit, bytesPerRuneIdx int) int {
	multiplier := float64(maxIdx) / float64(throwawayLimit)
	increasedLength := math.Ceil(multiplier * float64(length))
	return int(increasedLength) * bytesPerRuneIdx
}
