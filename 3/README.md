Zadanie na 3. zajęcia
=====================

Celem dzisiejszych zajęć jest przetestowanie funkcji, które służą
do wyszukiwania wzorca w tekście. Napiszą państwo testy jednostkowe,
dzięki którym sprawdza się, czy dana funkcja działa poprawnie, oraz
testy wydajnościowe, dzięki którym można porównać czas działania
różnych funkcji, których wyniki mają być jednakowe.

Testy jednostkowe
-----------------

1. Piszę testy jednostkowe w takich funkcjach:
`func TestFoo(t *testing.T)`

2. Proszę przeczytać sekcję dokumentacji o testach jednostkowych:
https://pkg.go.dev/testing (do sekcji *Benchmarks*)

3. Proszę w pliku `matching_test.go` napisać testy jednostkowe
tych funkcji:

* `Naive`
* `BackwardNaive`
* `BoyerMoore`
* `KMP`
* `KarpRabin`
* `ShiftOr`

W każdym teście należy:

* przypisać do zmiennych `pat` i `text` pewne wartości typu `[]byte`
* przepisać poniższy fragment kodu, odpowiednio zmieniając nazwę
wywoływanej funkcji:

```
got := []int{}
Naive(pat, text, func(n int) { got = append(got, n) })
want := indices(pat, text)
if !slices.Equal(got, want) {
	// Zgłoś błąd, korzystając z funkcji `t.Errorf`
}
```

4. Proszę uruchomić testy jednostkowe poleceniem `go test`

Testy wydajnościowe
-------------------

1. Piszę testy wydajnościowe w takich funkcjach:
`func BenchmarkFoo(b *testing.B)`

2. Proszę przeczytać sekcję dokumentacji o testach wydajnościowych:
https://pkg.go.dev/testing#hdr-Benchmarks (do opisu funkcji
`RunParallel`)

3. Proszę pobrać z serwisu https://wolnelektury.pl dowolną książkę
w formacie `.txt`

4. Proszę wybrać dwa wyrazy występujące w tej książce: jeden
4-literowy lub 5-literowy, a drugi co najmniej 10-literowy

5. Proszę w pliku `matching_test.go` napisać po 2 testy wydajnościowe
tych funkcji:

* `Naive`
* `BackwardNaive`
* `BoyerMoore`
* `KMP`
* `KarpRabin`
* `ShiftOr`

Jeden test ma mierzyć czas wyszukiwania krótszego wyrazu, a drugi —
czas wyszukiwania dłuższego wyrazu. Funkcje, które zawierają te testy,
można nazwać na przykład `BenchmarkShortNaive`, `BenchmarkLongNaive`
i tak dalej.

W każdym teście należy:

* wczytać do zmiennej `text` zawartość pobranej książki, korzystając
z funkcji `os.ReadFile` (https://pkg.go.dev/os#ReadFile)
* przypisać do zmiennej `pat` jeden z wybranych wyrazów jako wartość
typu `[]byte`
* wyzerować stoper, wywołując funkcję `b.ResetTimer`
* przepisać poniższy fragment kodu, odpowiednio zmieniając nazwę
wywoływanej funkcji:

```
for i := 0; i < b.N; i++ {
	Naive(pat, text, func(int) {})
}
```

6. Proszę uruchomić testy wydajnościowe poleceniem `go test -bench=.`

7. Proszę przepisać czas działania testowanych funkcji do arkusza
http://tiny.cc/at-lab3-2024
Proszę zaokrąglić czas działania tych funkcji do mikrosekund.
Na przykład jeśli otrzymam wyniki

```
BenchmarkShortNaive-8   	    1484	    779364 ns/op
```
to wpisuję do odpowiedniej komórki arkusza wartość 779
