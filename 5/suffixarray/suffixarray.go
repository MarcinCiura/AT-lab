package suffixarray

import (
	"bytes"
	"sort"
)

type suffix struct {
	s []byte
	n int
}

// sortedSuffixes zwraca wycinek, którego elementami są wszystkie
// sufiksy łańcucha `data`. Jeśli pewien sufiks s1 poprzedza inny
// sufiks s2 w porządku leksykograficznym, to indeks pozycji, od
// której zaczyna się sufiks s1, występuje w zwracanym wycinku
// przed indeksem pozycji, od której zaczyna się sufiks s2
func sortedSuffixes(data []byte) []suffix {
	suffixes := make([]suffix, len(data))
	for i := range suffixes {
		suffixes[i].s = data[i:]
		suffixes[i].n = i
	}
	sort.Slice(suffixes, func(i, j int) bool {
		return bytes.Compare(suffixes[i].s, suffixes[j].s) < 0
	})
	return suffixes
}

// Index implementuje tablicę sufiksów
type Index struct {
	data []byte
	ints []int
}

// New zwraca nowy Index dla łańcucha `data`
func New(data []byte) Index {
	suffixes := sortedSuffixes(data)
	index := Index{
		make([]byte, len(data)),
		make([]int, len(data)),
	}
	index.data = data
	for i := range suffixes {
		index.ints[i] = suffixes[i].n
	}
	return index
}

// Suffix zwraca ten sufiks łańcucha indeksowanego przez
// index `x`, który występuje na i-tej pozycji w tablicy
// sufiksów zawartej w indeksie `x`
func (x Index) Suffix(i int) []byte {
	return x.data[i:]
}

// LookupAll zwraca wycinek, którego elementami są wszystkie
// te pozycje, na których łańcuch `s` występuje w łańcuchu
// indeksowanym przez Index
func (x Index) LookupAll(s []byte) []int {
	i := sort.Search(len(x.ints), func(i int) bool {
		return bytes.Compare(x.Suffix(x.ints[i]), s) >= 0
	})
	j := i + sort.Search(len(x.ints)-i, func(j int) bool {
		return !bytes.HasPrefix(x.Suffix(x.ints[i+j]), s)
	})
	return x.ints[i:j]
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

// LongestRepeatingSubstring znajduje najdłuższy taki łańcuch,
// który występuje w łańcuchu `text` co najmniej `k` razy
func LongestRepeatingSubstring(text []byte, k int) []byte {
	index := New(text)
	substr := []byte{}
	for i := 0; i+k <= len(text); i++ {
		m := lenOfCommonPrefix(
			index.Suffix(i), index.Suffix(i+k-1))
		if m > len(substr) {
			substr = index.Suffix(i)[:m]
		}
	}
	return substr
}
