package nilsimsa

import (
	"testing"
)

func TestNilsimsa(t *testing.T) {
	s1 := "To niedźwiedź czy może dźwiedź? Chyba nie dźwiedź."
	s2 := "Czy to dźwiedź, czy niedźwiedź? Może nie dźwiedź."
	s3 := "Najgłupsze zwierzę w dżungli? Niedźwiedź polarny."
	data := []struct{
		s1, s2 string
		want   int
	}{
		{s1, s2, 47},
		{s1, s3, 82},
		{s2, s3, 83},
	}
	for _, d := range data {
		if got := HammingDistance(Nilsimsa(d.s1), Nilsimsa(d.s2)); got != d.want {
			t.Errorf("HammingDistance(Nilsimsa(%v), Nilsimsa(%v)) == %d, want %d",
				d.s1, d.s2, got, d.want)
		}
	}
}
