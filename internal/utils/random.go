package utils

import (
	"math/rand"
)

func CheckChance(chance float64) bool {
	return rand.Float64() < chance
}
