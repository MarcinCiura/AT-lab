package matching

import (
	"slices"
)

// findLastOccurrences zwraca mapę, która:
// + odwzorowuje wszystkie takie znaki, które występują
//   w łańcuchu `s` na `k+1`, gdzie `k` to indeks ostatniego
//   wystąpienia danego znaku w `s`
// + odwzorowuje wszystkie takie znaki, które nie występują
//   w łańcuchu `s`, na 0
func findLastOccurrences(s []byte) []int {
	lastOccurrences := make([]int, 256)
	for k, c := range s {
		lastOccurrences[c] = k + 1
	}
	return lastOccurrences
}

func simpleFindShiftOfSuffix(s []byte, i int) int {
	for k := 1; k <= i; k++ {
		if s[i-k] != s[i] &&
			slices.Equal(s[i+1:], s[i-k+1:len(s)-k]) {
			return k
		}
	}
	for k := i + 1; k < len(s); k++ {
		if slices.Equal(s[k:], s[:len(s)-k]) {
			return k
		}
	}
	return len(s)
}

func simpleComputeGoodSuffixes(s []byte) []int {
	r := make([]int, len(s))
	for i := range s {
		r[i] = simpleFindShiftOfSuffix(s, i)
	}
	return r
}

// boyerMooreHasPrefix zwraca parę (R, S). R ma wartość `true`,
// jeśli `slices.Equal(text[:len(pat)], pat)`; S określa,
// o ile pozycji w prawo należy przesunąć wzorzec `pat`
func boyerMooreHasPrefix(text, pat []byte,
	lastOccurrences []int, goodSuffixes []int) (bool, int) {
	for i := len(pat) - 1; i >= 0; i-- {
		if text[i] != pat[i] {
			return false, max(i+1 - lastOccurrences[text[i]],
				goodSuffixes[i])
		}
	}
	return true, goodSuffixes[0]
}

// BoyerMoore wywołuje `output(i)` dla każdego takiego `i`,
// że `slices.Equal(text[i:i+len(pat)], pat)`
func BoyerMoore(pat, text []byte, output func(int)) {
	lastOccurrences := findLastOccurrences(pat)
	goodSuffixes := simpleComputeGoodSuffixes(pat)
	for i := 0; i+len(pat) <= len(text); /**/ {
		found, shift := boyerMooreHasPrefix(text[i:], pat,
			lastOccurrences, goodSuffixes)
		if found {
			output(i)
		}
		i += shift
	}
}

// setNthBit zwraca maskę, w której bit na pozycji n jest równy 1
func setNthBit(n int) uint64 {
	return uint64(1) << n
}

// nthBit zwraca bit na pozycji n w masce m
func nthBit(m uint64, n int) uint64 {
	return (m >> n) & 1
}

// makeMask zwraca tablicę 256 masek; bity i-tej maski są równe 0
// na pozycjach równych wszystkim pozycjom znaku i we wzorcu pat
func makeMask(pat []byte) [256]uint64 {
	m := [256]uint64{}
	for c := 0; c < 256; c++ {
		m[c] = ^uint64(0) // Ustaw wszystkie bity maski m[c]
	}
	for j, c := range pat {
		m[c] &^= setNthBit(j) // Wyzeruj j-ty bit maski m[c]
	}
	// Dla 0 <= c < 256, 0 <= j < len(pat) zachodzi
	// (nthBit(m[c], j) == 0) == (pat[j] == c)
	return m
}

// FuzzyShiftOrH wywołuje funkcję `output(i)` dla każdego
// takiego indeksu `i`, że `text[i:i+len(pat)]` różni się
// od `pat` co najwyżej na 2 pozycjach
func FuzzyShiftOrH(pat, text []byte, output func(int)) {
	m := makeMask(pat)
	s0, s1, s2 := ^uint64(0), ^uint64(0), ^uint64(0)
	for i, c := range text {
		// Uwzględnij zamianę 1 znaku
		s2 = ((s2 << 1) | m[c]) & (s1 << 1)
		s1 = ((s1 << 1) | m[c]) & (s0 << 1)
		s0 = (s0 << 1) | m[c]
		if nthBit(s2, len(pat)-1) == 0 {
			output(i - len(pat) + 1)
		}
	}
}

// FuzzyShiftOrL wywołuje funkcję `output(i)` dla każdego takiego
// indeksu `i`, że odległość Levenshteina między pewnym wycinkiem
// `text[...:i+1]` a wzorcem `pat` wynosi co najwyżej 2
func FuzzyShiftOrL(pat, text []byte, output func(int)) {
	m := makeMask(pat)
	s0, s1, s2 := ^uint64(0), ^uint64(0), ^uint64(0)
	for i, c := range text {
		// Uwzględnij zamianę 1 znaku lub wstawienie 1 znaku
		s2 = ((s2 << 1) | m[c]) & (s1 << 1) & s1
		s1 = ((s1 << 1) | m[c]) & (s0 << 1) & s0
		s0 = (s0 << 1) | m[c]
		// Uwzględnij usunięcie 1 znaku
		s1 &= (s0 << 1)
		s2 &= (s1 << 1)
		if nthBit(s2, len(pat)-1) == 0 {
			output(i) // Zwróć pozycję ostatniego znaku wycinka
		}
	}
}
