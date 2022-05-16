package utils

import (
	"math/rand"
	"time"
)

func CheckChance(chance float64) bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64() < chance
}
