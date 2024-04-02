package matching

import (
	"slices"
	"testing"
)

func TestPreprocess(t *testing.T) {
	data := []string{
		"aaaaaaa",
		"pies",
		"dźwiedź",
		"owocowo",
		"indianin",
		"nienapełnienie",
	}
	for _, in := range data {
		got := Preprocess([]byte(in))
		want := SimplePreprocess([]byte(in))
		if !slices.Equal(got, want) {
			t.Errorf(`Preprocess(%#v) == %#v want %#v`,
				in, got, want)
		}
	}
}

func indices(pat, text []byte) []int {
	r := []int{}
	for i := 0; i+len(pat) <= len(text); i++ {
		if slices.Equal(text[i:i+len(pat)], pat) {
			r = append(r, i)
		}
	}
	return r
}
