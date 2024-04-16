Zadanie na 4. zajęcia
=====================

Celem dzisiejszych zajęć jest:

* Zapoznać się z "przybliżonymi" wersjami algorytmu Shift-Or. Te
  wersje algorytmu Shift-Or znajdują takie podłańcuchy tekstu, których
  odległość Hamminga lub odległość Levenshteina od wzorca nie
  przekracza danej stałej, w naszym przypadku stałej 2.

* Przetestować pakiet, który służy do znajdowania wielu wzorców w
  tekście algorytmem Aho-Corasick. Napiszą państwo testy jednostkowe,
  które porównają wynik działania algorytmu Aho-Corasick z wynikiem
  działania wielokrotnie uruchomionego algorytmu Boyera-Moore'a, oraz
  testy wydajnościowe, które porównają czas działania algorytmu
  Aho-Corasick z czasem działania wielokrotnie uruchomionego algorytmu
  Boyera-Moore'a.

1. Proszę wykonać polecenie

```
go get github.com/BobuSumisu/aho-corasick
```

2. Proszę pobrać z serwisu https://wolnelektury.pl dowolną książkę
w formacie `.txt`

Wersje algorytmu Shift-Or
-------------------------

3. Proszę napisać funkcje `TestFuzzyShiftOrH` i `TestFuzzyShiftOrL`. W
obu tych funkcjach proszę kolejno:

* wczytać do zmiennej `text` zawartość pobranej książki, korzystając z
  funkcji `os.ReadFile` (https://pkg.go.dev/os#ReadFile)

* przypisać do zmiennej `pat` jakiś 5- lub 6-literowy wyraz jako
  wartość typu `[]byte`, czyli na przykład `pat := []byte("domek")`

* przypisać do zmiennej `got` pusty wycinek złożony z łańcuchów:
  `got := []string{}`

* wywołać odpowiednio funkcję `FuzzyShiftOrH` lub `FuzzyShiftOrL`,
  podając jako jej trzeci argument funkcję anonimową odpowiednio
  `func(n int) { got = append(got, string(text[n:n+len(pat)])) }` i
  `func(n int) { got = append(got, string(text[n-len(pat)-1:n+1])) }`;
  ciała tych funkcji anonimowych są różne, ponieważ funkcja
  `FuzzyShiftOrH` zwraca indeks początku wystąpienia wzorca w tekście,
  a funkcja `FuzzyShiftOrL` zwraca indeks końca wystąpienia wzorca w
  tekście

* wypisać zawartość wycinka `got`, korzystając z funkcji `fmt.Printf`
  i specyfikatora formatu `%v`; niech funkcje `TestFuzzyShiftOrH` i
  `TestFuzzyShiftOrL` nie robią żadnych testów, tylko wypisują
  zawartość wycinka `got`

* wykonać polecenie `go test` i zobaczyć wyniki działania tych dwóch
  funkcji

* przepisać kilka znalezionych przykładów do arkusza
  http://tiny.cc/at-lab4-2024

Algorytm Aho-Corasick
---------------------

4. Proszę wypełnić zmienną globalną `words` wszystkimi formami odmiany
wybranego rzeczownika lub przymiotnika.

5. Proszę napisać funkcję `TestAhoCorasick`. W tej funkcji proszę
kolejno:

* wczytać do zmiennej `text` zawartość pobranej książki, korzystając z
  funkcji `os.ReadFile` (https://pkg.go.dev/os#ReadFile)

* zbudować drzewo trie, z którego korzysta algorytm Aho-Corasick:

```
builder := ahocorasick.NewTrieBuilder()
builder.AddStrings(words)
trie := builder.Build()
```

* wyszukać algorytmem Aho-Corasick wszystkie wystąpienia w tekście
  wszystkich wzorców zapisanych w zmiennej `words` i przypisać je do
  zmiennej `matches`:

```
matches := trie.MatchString(string(text))
```

* napisać część funkcji, która wypełnia n-ty element tablicy wycinków
  `got` numerami tych pozycji w tekście, na których zaczynają się
  wystąpienia n-tego elementu wycinka `words`, znalezione algorytmem
  Aho-Corasick; proszę zmienić wartość stałej 12 na liczbę równą
  długości wycinka `words`:

```
var got [12][]int64
for _, m := range matches {
	got[m.Pattern()] = append(got[m.Pattern()], m.Pos())
}
```

* napisać część funkcji, która która wypełnia n-ty element tablicy
  wycinków `want` numerami tych pozycji w tekście, na których
  zaczynają się wystąpienia n-tego elementu wycinka `words`,
  znalezione algorytmem Boyera-Moore'a; proszę zmienić wartość stałej
  12 na liczbę równą długości wycinka `words`:

```
var want [12][]int64
for i, pat := range words {
	BoyerMoore([]byte(pat), text, func(n int) {
		want[i] = append(want[i], int64(n))
	})
}
```

* napisać część funkcji, która porównuje tablicę wycinków `got` z
  tablicą wycinków `want`:

```
for i := range got {
	if !slices.Equal(got[i], want[i]) {
		t.Errorf("got[%d] == %v want %v", i, got[i], want[i])
	}
}
```

* dla zaspokojenia ciekawości można dopisać do tej funkcji fragment:

```
fmt.Printf("Znalezione wystąpienia wyrazów %v:\n%v\n", words, got)
```

6. Proszę napisać funkcję `BenchmarkAhoCorasick`. W tej funkcji proszę
kolejno:

* wczytać do zmiennej `text` zawartość pobranej książki, korzystając z
  funkcji `os.ReadFile` (https://pkg.go.dev/os#ReadFile)

* przypisać do zmiennej `stext` zawartość zmiennej `text` jako łańcuch:

```
stext := string(text)
```

* wyzerować stoper, wywołując funkcję `b.ResetTimer`

* przepisać poniższy fragment kodu, który mierzy czas budowania drzewa
  trie z wzorców i wyszukiwania tych wzorców w tekście:

```
for i := 0; i < b.N; i++ {
	builder := ahocorasick.NewTrieBuilder()
	builder.AddStrings(words)
	trie := builder.Build()
	trie.MatchString(stext)
}
```

7. Proszę napisać funkcję `BenchmarkBoyerMoore`, żeby móc porównywać
   czas wyszukiwania wszystkich wzorców algorytmem Aho-Corasick i
   wielokrotnie uruchamianym algorytmem Boyera-Moore'a. Funkcja
   `BenchmarkBoyerMoore` jest podobna do funkcji
   `BenchmarkAhoCorasick` z tym, że wyszukiwanie wszystkich wzorców
   w tekście jest zakodowane tak:

```
for _, pat := range words {
	BoyerMoore([]byte(pat), text, func(int){})
}
```

8. Proszę porównać czas wyszukiwania wszystkich wzorców algorytmem
Aho-Corasick i algorytmem Boyera-Moore'a i wpisać go do arkusza
http://tiny.cc/at-lab4-2024 — przypominam, że testy wydajnościowe
uruchamia się poleceniem `go test -bench=.`
