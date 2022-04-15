package bot

import (
	"math/rand"
	"time"
)

func checkChance(chance float64) bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64() < chance
}
