package bot

import (
	"math/rand"
	"testing"
)

func Test_getRandomPeka(t *testing.T) {
	rand.Seed(1)
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ok",
			want: ":uuu:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRandomEmoji(); got != tt.want {
				t.Errorf("getRandomPeka() = %v, want %v", got, tt.want)
			}
		})
	}
}
