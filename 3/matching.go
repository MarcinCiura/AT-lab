package matching

import (
	"slices"
)

// hasPrefix zwraca `true`, jeśli
// `slices.Equal(s[:len(pat)], pat)`
func hasPrefix(s, pat []byte) bool {
	for j := 0; j < len(pat); j++ {
		if s[j] != pat[j] {
			return false
		}
	}
	return true
}

// Naive wywołuje funkcję `output(i)` dla każdego takiego
// indeksu `i`, że `slices.Equal(text[i:i+len(pat)], pat)`
func Naive(pat, text []byte, output func(int)) {
	for i := 0; i+len(pat) <= len(text); i++ {
		if hasPrefix(text[i:], pat) {
			output(i)
		}
	}
}

// backwardHasPrefix zwraca `true`, jeśli
// `slices.Equal(s[:len(pat)], pat)`
func backwardHasPrefix(s, pat []byte) bool {
	for i := len(pat) - 1; i >= 0; i-- {
		if s[i] != pat[i] {
			return false
		}
	}
	return true
}

// BackwardNaive wywołuje `output(i)` dla każdego takiego `i`,
// że `slices.Equal(text[i:i+len(pat)], pat)`
func BackwardNaive(pat, text []byte, output func(int)) {
	for i := 0; i+len(pat) <= len(text); i++ {
		if backwardHasPrefix(text[i:], pat) {
			output(i)
		}
	}
}

// lenOfCommonPrefix zwraca długość najdłuższego
// wspólnego prefiksu łańcuchów `s` i `t`
func lenOfCommonPrefix(s, t []byte) int {
	k := 0
	for ; k < min(len(s), len(t)); k++ {
		if s[k] != t[k] {
			return k
		}
	}
	return k
}

// SimplePreprocess zwraca wycinek. `k`-ty element tego
// wycinka jest równy długości najdłuższego takiego prefiksu
// łańcucha `s[k:]`, który jest równy pewnemu prefiksowi
// łańcucha `s`
func SimplePreprocess(s []byte) []int {
	z := make([]int, len(s))
	for k := 1; k < len(s); k++ {
		z[k] = lenOfCommonPrefix(s[k:], s)
	}
	return z
}

// Preprocess zwraca wycinek. `k`-ty element tego wycinka
// jest równy długości najdłuższego takiego prefiksu
// łańcucha `s[k:]`, który jest równy pewnemu prefiksowi
// łańcucha `s`
func Preprocess(s []byte) []int {
	z := make([]int, len(s))
	l := 0
	r := 0
	for k := 1; k < len(s); k++ {
		if k >= r {
			z[k] = lenOfCommonPrefix(s[k:], s)
			if z[k] > 0 {
				l = k
				r = k + z[k]
			}
		} else if z[k-l] >= r-k {
			z[k] = r - k + lenOfCommonPrefix(s[r:], s[r-k:])
			l = k
			r = k + z[k]
		} else {
			z[k] = z[k-l]
		}
		// bytes.HasPrefix(s, s[k:k+z[k]])
		// bytes.HasPrefix(s, s[l:r])
	}
	return z
}

// findLastOccurrences zwraca mapę, która:
//   - odwzorowuje wszystkie takie znaki, które występują
//     w łańcuchu `s` na `k+1`, gdzie `k` to indeks ostatniego
//     wystąpienia danego znaku w `s`
//   - odwzorowuje wszystkie takie znaki, które nie występują
//     w łańcuchu `s`, na 0
func findLastOccurrences(s []byte) map[byte]int {
	lastOccurrences := map[byte]int{}
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
	lastOccurrences map[byte]int, goodSuffixes []int) (bool, int) {
	for i := len(pat) - 1; i >= 0; i-- {
		if text[i] != pat[i] {
			return false, max(i+1-lastOccurrences[text[i]],
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

// KMPPrefixFunction zwraca wycinek. `j`-ty element tego wycinka
// to wartość funkcji prefiksowej `p[j]`, czyli długość najdłuższego
// takiego właściwego sufiksu łańcucha `s[:j+1]`, który jest pewnym
// prefiksem łańcucha `s`
func KMPPrefixFunction(s []byte) []int {
	z := Preprocess(s)
	p := make([]int, len(s)+1)
	for j := len(s) - 1; j > 0; j-- {
		p[j+z[j]] = z[j]
	}
	return p
}

func KMP(pat, text []byte, output func(int)) {
	// len(pat) > 0
	p := KMPPrefixFunction(pat)
	j := 0
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pat[j] {
			j = p[j]
		}
		if text[i] == pat[j] {
			j++
		}
		if j == len(pat) {
			output(i - len(pat) + 1)
			j = p[j]
		}
	}
}

// hashByteModN zwraca liczbę całkowitą z przedziału [0, n)
func hashByteModN(b byte, h, n uint64) uint64 {
	// Nie ma przepełnienia, jeśli
	// h<<8 + uint64(b) < 1<<64
	// h<<8 < 1<<64 - uint64(b)
	// h<<8 < 1<<64 - 1<<8 + 1
	// h<<8 <= 1<<64 - 1<<8
	// h<<8>>8 <= 1<<64>>8 - 1<<8>>8
	// h <= 1<<56 - 1
	// h < 1<<56
	return (h<<8 + uint64(b)) % n
}

// hashBytesModN zwraca liczbę całkowitą z przedziału [0, n)
func hashBytesModN(bs []byte, n uint64) uint64 {
	h := uint64(0)
	for _, b := range bs {
		// h < n
		h = hashByteModN(b, h, n)
	}
	return h
}

// twoToPower8PModN zwraca 2**(8*p) % n
func twoToPower8PModN(p int, n uint64) uint64 {
	r := uint64(1)
	for i := 0; i < p; i++ {
		// r == 1<<(8*i) % n
		// r == 2**(8*i) % n
		// Nie ma przepełnienia, jeśli
		// r<<8 < 1<<64
		// r<<8>>8 < 1<<64>>8
		// r < 1<<56
		r = (r << 8) % n
	}
	return r
}

// unhashByteModN zwraca liczbę całkowitą z przedziału [0, n)
func unhashByteModN(b byte, h, n, power uint64) uint64 {
	// Nie ma przepełnienia, jeśli
	// h - power*uint64(b) > -(1<<64)
	// 0 - power*(1<<8-1) > -(1<<64)
	// power*(1<<8-1) < 1<<64
	// power<<8 - power < 1<<64
	// power<<8 < 1<<64
	// power<<8>>8 < 1<<64>>8
	// power < 1<<56
	r := h - power*uint64(b)
	if int64(r) < 0 {
		r += n * (-r/n + 1)
	}
	return r
}

// Największa liczba pierwsza mniejsza niż 1<<56
// https://t5k.org/lists/2small/0bit.html
const N uint64 = 1<<56 - 5

func KarpRabin(pat, text []byte, output func(int)) {
	// len(pat) <= len(text)
	h := hashBytesModN(text[:len(pat)], N)
	ph := hashBytesModN(pat, N)
	power := twoToPower8PModN(len(pat), N)
	for i := 0; ; i++ {
		// h == hashBytesModN(text[i:i+len(pat)], N)
		if h == ph && slices.Equal(pat, text[i:i+len(pat)]) {
			output(i)
		}
		if i+len(pat) >= len(text) {
			break
		}
		h = hashByteModN(text[i+len(pat)], h, N)
		h = unhashByteModN(text[i], h, N, power)
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

func ShiftOr(pat, text []byte, output func(int)) {
	// len(pat) != 0
	m := makeMask(pat)
	s := ^uint64(0) // Ustaw wszystkie bity maski s
	for i, c := range text {
		// Dla 1 < j < min(len(pat), i) zachodzi
		// (nthBit(s, j-1) == 0) ==
		//    slices.Equal(pat[:j], text[i-j:i])
		s = (s << 1) | m[c] // Shift-Or
		if nthBit(s, len(pat)-1) == 0 {
			output(i - len(pat) + 1)
		}
	}
}
