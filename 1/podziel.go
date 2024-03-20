package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

// Wyraz i jego punkty
type Pair struct {
	word   string
	points int
}

// main kolejno:
// + czyta wyrazy z pliku "slowa.txt"
// + na wszystkie możliwe sposoby dzieli każdy wyraz na 2 części
// + jeśli obie części są wyrazami, przyznaje tym częściom po 1 punkcie
// + sortuje takie części, które mają więcej niż 0 punktów, według
//   malejącej kolejności punktów
// + wypisuje te części i ich punkty
// + wypisuje takie części, które mają tyle samo punktów, w porządku
//   leksykograficznym
func main() {
	wordList := ReadLines("slowa.txt")
	wordSet := map[string]bool{}
	for _, w := range wordList {
		wordSet[w] = true
	}
	wordCounter := map[string]int{}
	for _, w := range wordList {
		for _, parts := range Split(w) {
			IncrementIfBothIn(parts, wordSet, &wordCounter)
		}
	}
	pairs := Sort(wordCounter)
	for _, p := range pairs {
		fmt.Printf("%d %s\n", p.points, p.word)
	}
}

// Read wczytuje wiersze z pliku o nazwie `filename` i zwraca je w
// wycinku tablicy łańcuchów
func ReadLines(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}

// Split na wszystkie możliwe sposoby dzieli wyraz `word` na 2
// niepuste łańcuchy i zwraca te 2 łańcuchy w wycinku tablicy
// wycinków tablic łańcuchów
func Split(word string) [][]string {
	r := [][]string{}
	for i := 1; i < len(word); i++ {
		r = append(r, []string{word[:i], word[i:]})
	}
	return r
}

// IncrementIfBothIn zwiększa licznik `counter` przy tych łańcuchach z
// wycinka `parts`, które należą do zbioru łańcuchów `set`
func IncrementIfBothIn(
	parts []string, set map[string]bool, counter *map[string]int) {
	if set[parts[0]] && set[parts[1]] {
		(*counter)[parts[0]]++
		(*counter)[parts[1]]++
	}
}

// Sort sortuje wycinek par (word, points). Gdy 2 pary mają różną
// liczbę punktów, ta para, która ma więcej punktów, poprzedza tę
// parę, który ma mniej punktów. Gdy 2 pary mają tyle samo punktów, ta
// para, której pole `word` jest pierwsze w porządku
// leksykograficznym, poprzedza drugą parę
func Sort(counter map[string]int) []Pair {
	r := []Pair{}
	for w, p := range counter {
		r = append(r, Pair{w, p})
	}
	slices.SortFunc(r, func(a, b Pair) int {
		if n := b.points - a.points; n != 0 {
			return n
		}
		return cmp.Compare(a.word, b.word)
	})
	return r
}
