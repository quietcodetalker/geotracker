package testutil

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomFloat64 returns random number between given min and max.
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomInt returns a random number in the interval [min, max].
func RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}

	return min + rand.Intn(max-min+1)
}

const (
	// CharacterSetAlphabet is a character set of all lowercase and uppercase letters.
	CharacterSetAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// CharacterSetAlphanumeric is a character set of all characters from CharacterSetAlphabet plus digits from 0 to 9.
	CharacterSetAlphanumeric = CharacterSetAlphabet + "0123456789"
)

// RandomString returns a string minLength to maxLength characters long of random characters
// taken from characterSet.
func RandomString(minLength, maxLength int, characterSet string) string {
	if minLength < 0 {
		minLength = 0
	}
	if maxLength < 0 {
		maxLength = 0
	}

	var sb strings.Builder
	k := len(characterSet)
	n := RandomInt(minLength, maxLength)

	for i := 0; i < n; i++ {
		c := characterSet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

const (
	UsernameCharacterSet = CharacterSetAlphanumeric
	UsernameMinLen       = 4
	UsernameMaxLen       = 16
)

// RandomUsername returns a random string that suits predefined username conditions.
func RandomUsername() string {
	return RandomString(UsernameMinLen, UsernameMaxLen, UsernameCharacterSet)
}

// RandomTimeInterval returns two random timestamps where the first one is before the second one.
func RandomTimeInterval() (time.Time, time.Time) {
	hours := RandomInt(-24, 0)
	from := time.Now().Add(time.Hour * time.Duration(hours))
	to := from.Add(time.Second * time.Duration(RandomInt(1, 1000)))

	return from, to
}
