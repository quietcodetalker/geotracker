package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomFloat64 returns random number between given min and max.
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
