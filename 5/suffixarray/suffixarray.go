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
// sufiksy łańcucha `text`. Jeśli pewien sufiks S1 poprzedza inny
// sufiks S2 w porządku leksykograficznym, to indeks tej pozycji,
// od której zaczyna się sufiks S1, występuje w zwracanym wycinku
// przed indeksem tej pozycji, od której zaczyna się sufiks S2
func sortedSuffixes(text []byte) []suffix {
	suffixes := make([]suffix, len(text))
	for i := range suffixes {
		suffixes[i].s = text[i:]
		suffixes[i].n = i
	}
	sort.Slice(suffixes, func(i, j int) bool {
		return bytes.Compare(suffixes[i].s, suffixes[j].s) < 0
	})
	return suffixes
}

// Index implementuje indeks łańcucha `text`. Ten indeks
// składa się z tablicy sufiksów łańcucha `text`
// i z łańcucha `text`
type Index struct {
	suffixes []int
	text     []byte
}

// New zwraca indeks łańcucha `text`
func New(text []byte) Index {
	index := Index{
		make([]int, len(text)),
		make([]byte, len(text)),
	}
	index.text = text
	for i, suf := range sortedSuffixes(text) {
		index.suffixes[i] = suf.n
	}
	return index
}

// Suffix zwraca ten sufiks łańcucha `x.text`, który
// zaczyna się od `i`-tej pozycji tego łańcucha
func (x Index) Suffix(i int) []byte {
	return x.text[i:]
}

// LookupAll zwraca wycinek, którego elementami są wszystkie
// te pozycje, na których łańcuch `s` występuje w łańcuchu `x.text`
func (x Index) LookupAll(s []byte) []int {
	i := sort.Search(len(x.suffixes), func(i int) bool {
		return bytes.Compare(x.Suffix(x.suffixes[i]), s) >= 0
	})
	j := i + sort.Search(len(x.suffixes)-i, func(j int) bool {
		return !bytes.HasPrefix(x.Suffix(x.suffixes[i+j]), s)
	})
	return x.suffixes[i:j]
}
