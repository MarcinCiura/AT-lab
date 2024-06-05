# Zadanie 1.

1. Proszę zmienić sygnaturę funkcję `FuzzyShiftOrL` z 4. zajęć na taką:

```go
// FuzzyShiftOrL wywołuje funkcję `output(i, j)` dla każdej takiej
// pary indeksów (`i`, `j`), że odległość Levenshteina między
// wycinkiem `text[i:j]` a wzorcem `pat` wynosi co najwyżej 2
func FuzzyShiftOrL(pat, text []byte, output func(int, int)) {
```

2. Proszę zmienić ciało funkcji `FuzzyShiftOrL` tak, aby jej działanie
było zgodne z jej sygnaturą. Mogą państwo skorzystać z algorytmu
Wagnera-Fischera lub z algorytmu Allisona, aby znajdować taki wycinek
`text[...:j]`, którego odległość edycyjna od wzorca `pat` wynosi
co najwyżej 2. Kody funkcji, które implementują oba te algorytmy,
są przedstawione na slajdach do wykładu 6.

3. Proszę pamiętać o testach jednostkowych

4. Proszę przepisać do arkusza http://tiny.cc/at-lab6-2024 kilka
przykładów działania zmienionej funkcji `FuzzyShiftOrL` na dowolnej
książce pobranej z serwisu https://wolnelektury.pl


# Zadanie 2.

1. Proszę pobrać z dowolnego serwisu informacyjnego I1 treść dowolnej
wiadomości z bieżącego dnia W1 i zapisać ją jako tekst

2. Proszę pobrać z innego serwisu informacyjnego I2 treści 5 dowolnych
wiadomości W2-W6 z bieżącego dnia i zapisać je jako teksty. Jedna z
tych wiadomości powinna dotyczyć tego samego zdarzenia, co wiadomość
W1 pobrana z serwisu I1

3. Proszę wyznaczyć funkcję haszującą Nilsimsa 6 pobranych wiadomości

4. Proszę wyznaczyć odległość Hamminga między parami wyników funkcji
Nilsimsa wyznaczonych dla wiadomości W1-W2, W1-W3,..., W1-W6

5. Proszę przepisać wyniki zadania do arkusza
http://tiny.cc/at-lab6-2024

Algorytm Nilsimsa jest omówiony na slajdach do wykładu 6 i
zaimplementowany w module `github.com/MarcinCiura/AT-lab/6/nilsimsa`