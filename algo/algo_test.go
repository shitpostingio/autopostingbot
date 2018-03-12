package algo

import (
	"testing"
)

func TestAlgo(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int64
	}{
		{
			"base test, no elements",
			[]string{},
			0,
		},
		{
			"exactly twelve elements in the array",
			generaStringa(12),
			1,
		},
		{
			"exactly 24 elements in the array",
			generaStringa(24),
			1,
		},
		{
			"exactly 26 elements in the array",
			generaStringa(26),
			2,
		},
		{
			"exactly 48 elements in the array",
			generaStringa(48),
			2,
		},
		{
			"exactly 50 elements in the array",
			generaStringa(50),
			2,
		},
		{
			"exactly 72 elements in the array",
			generaStringa(72),
			3,
		},
		{
			"exactly 130 elements in the array",
			generaStringa(130),
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Algo(tt.args); got != tt.want {
				t.Errorf("Algo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generaStringa(lunghezza int) []string {
	var s []string
	for i := 0; i < lunghezza; i++ {
		s = append(s, "i")
	}
	return s
}
