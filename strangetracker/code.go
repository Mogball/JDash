package strangetracker

import (
	"strings"
	"strconv"
	"log"
	"math"
)

type CodeSet struct {
	A     int64
	B     int64
	C     int64
	Mod   int64
	Limit int64
}

func splitAndGetBytes(encodedString string, separator string) []byte {
	pieces := strings.Split(encodedString, separator)
	bytes := make([]byte, len(pieces))
	for i, _ := range pieces {
		ordinal, err := strconv.Atoi(pieces[i])
		if err != nil {
			log.Println(err)
			bytes[i] = 0
			continue
		}
		bytes[i] = byte(ordinal)
	}
	return bytes
}

func computeNextOffset(code *CodeSet, i int64) int64 {
	return (code.C + code.B*i + code.A*i*i) % code.Mod
}

func encodeBytes(bytes []byte, code *CodeSet) []byte {
	result := make([]byte, len(bytes))
	var i int64 = 0
	for j, _ := range bytes {
		i = computeNextOffset(code, i)
		result[j] = byte((int64(bytes[j])*i + i) % code.Limit)
	}
	return result
}

func computeMaxDivFactor(code *CodeSet) int64 {
	maxRaw := code.Limit*code.Mod + code.Mod
	maxDivFactor := float64(maxRaw) / float64(code.Limit)
	return int64(math.Ceil(maxDivFactor))
}

func joinCode(encodedBytes []byte, separator string) string {
	pieces := make([]string, len(encodedBytes))
	for i, _ := range encodedBytes {
		pieces[i] = strconv.Itoa(int(encodedBytes[i]))
	}
	return strings.Join(pieces, separator)
}

func EncodeString(source string, separator string, code *CodeSet) string {
	encodedBytes := EncodeStringAsBytes(source, code)
	return joinCode(encodedBytes, separator)
}

func EncodeStringAsBytes(source string, code *CodeSet) []byte {
	bytes := []byte(source)
	return encodeBytes(bytes, code)
}

func crackOnce(chr byte, i int64, maxL int64, code *CodeSet) byte {
	for l := int64(0); l <= maxL; l++ {
		factor := float64((int64(chr)+l*code.Limit)-i) / float64(i)
		if factor == math.Floor(factor) {
			return byte(factor)
		}
	}
	return 0
}

func CrackBytes(encoded []byte, code *CodeSet) []byte {
	var i int64 = 0
	maxL := computeMaxDivFactor(code)
	result := make([]byte, len(encoded))
	for j, chr := range encoded {
		i = computeNextOffset(code, i)
		result[j] = crackOnce(chr, i, maxL, code)
	}
	return result
}

func CrackString(encodedString string, separator string, code *CodeSet) string {
	encodedBytes := splitAndGetBytes(encodedString, separator)
	crackedBytes := CrackBytes(encodedBytes, code)
	return string(crackedBytes)
}
