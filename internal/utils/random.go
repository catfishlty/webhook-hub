package utils

import (
	"math/rand"
	"sync"
	"time"
)

type Random struct {
	r rand.Rand
}

var randomInstance *Random
var randomOnce sync.Once
var lowerLetterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var upperLetterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("01234567890")

func NewRandom() *Random {
	randomOnce.Do(func() {
		randomInstance = &Random{
			r: *rand.New(rand.NewSource(time.Now().UnixNano())),
		}
	})
	return randomInstance
}

func (r *Random) String(n int) string {
	return stringRunes(append(append(lowerLetterRunes, upperLetterRunes...), numberRunes...), n)
}

func (r *Random) StringWithLowerLetters(n int) string {
	return stringRunes(lowerLetterRunes, n)
}

func (r *Random) StringWithUpperLetters(n int) string {
	return stringRunes(upperLetterRunes, n)
}

func (r *Random) StringWithNumberLetters(n int) string {
	return stringRunes(numberRunes, n)
}

func (r *Random) StringWithLetters(n int) string {
	return stringRunes(append(lowerLetterRunes, upperLetterRunes...), n)
}

func (r *Random) StringWithLowerNumberLetters(n int) string {
	return stringRunes(append(lowerLetterRunes, numberRunes...), n)
}

func (r *Random) StringWithUpperNumberLetters(n int) string {
	return stringRunes(append(upperLetterRunes, numberRunes...), n)
}

func stringRunes(letterRunes []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
