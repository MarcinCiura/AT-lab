package main

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"regexp"
	"strings"
	"unicode"
)

// W tym pliku jest 5 zadań. Żeby rozwiązać każde z tych zadań,
// proszę napisać właściwy regexp i poprawnie go użyć.

var (
	// Pasuje do takich numerów telefonów stacjonarnych, które
	// zaczynają się od numeru kierunkowego 12 lub 91. Tylko
	// Ośrodek Wczasowy Łukęcin w województwie zachodniopomorskim
	// ma numer stacjonarny w innej strefie niż 12
	phone2322 = regexp.MustCompile(
		`(\+48 *)?(12|91)[ -]*(\d\d\d)[ -]*(\d\d)[ -]*(\d\d)`)
	// Pasuje do takich numerów telefonów komórkowych, których
	// cyfry są podawane po trzy
	phone333 = regexp.MustCompile(
		`(\+48 *)?(\d\d\d)[ -]*(\d\d\d)[ -]*(\d\d\d)`)
	// Pasuje do takich numerów telefonów komórkowych, których
	// cyfry są podawane po dwie
	phone3222 = regexp.MustCompile(``) // Zadanie 1
)

// TransformPhoneNumbers przekształca wszystkie numery telefonów w
// łańcuchu `s` na format 12-345-67-89 lub 500-123-456
func TransformPhoneNumbers(s string) string {
	s = phone2322.ReplaceAllString(s, "${2}-${3}-${4}-${5}")
	s = phone333.ReplaceAllString(s, "${2}-${3}-${4}")
	s = "" // Zadanie 1
	return s
}

// ToASCIIString zmienia każdy taki run łańcucha `s`, który jest
// literą alfabetu polskiego, na odpowiednią małą literę alfabetu
// łacińskiego
func ToASCIIString(s string) ASCIIString {
	s = strings.ToLower(s)
	t := transform.Chain(
		norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, err := transform.String(t, s)
	if err != nil {
		log.Fatal(err)
	}
	return ASCIIString(strings.ReplaceAll(s, "ł", "l"))
}

// Pasuje do ciągu dowolnych znaków, przed którym występuje 0 lub
// więcej znaków przestankowych, i po którym następuje 0 lub więcej
// znaków przestankowych
var reWord = regexp.MustCompile(`^\p{P}*(.*?)\p{P}*$`)

// ToASCIIWord usuwa z łańcucha `s` początkowe i końcowe ciągi 0 lub
// więcej znaków przestankowych
func ToASCIIWord(s ASCIIString) ASCIIWord {
	if m := reWord.FindStringSubmatch(string(s)); m != nil {
		return ASCIIWord(m[1])
	}
	return ASCIIWord(s)
}

var (
	// Pasuje do rzeczowników zakończonych na -cie, np. "Zygmuncie"
	reCieT = regexp.MustCompile("^(.+)cie$")
	// Pasuje do rzeczowników zakończonych na -dzie, np. "Alfredzie"
	reDzieD = regexp.MustCompile("^(.+)dzie$")
	// Pasuje do rzeczowników zakończonych na -ce, np. "Monice"
	reCeK = regexp.MustCompile("x") // Zadanie 2
	// Pasuje do przymiotników zakończonych na -cy, np. "Kowalscy"
	reCyK = regexp.MustCompile("^(.+)cy$")
	// Pasuje do rzeczowników zakończonych na -dze i -dzy,
	// np. "koledze" i "koledzy"
	reDzeG = regexp.MustCompile("x") // Zadanie 3
	// Pasuje do rzeczowników zakończonych na -rze i -rzy,
	// np. "Piotrze" i "doktorzy"
	reRzeR = regexp.MustCompile("x") // Zadanie 4
	// Pasuje do rzeczowników zakończonych na -cień,
	// np. "Kwiecień" i "Pierścień"
	reCień = regexp.MustCompile("x") // Zadanie 5
	// Pasuje do rzeczowników zakończonych na -dziec i -dzień,
	// np. "Dudziec", "Grudzień" i "Moździeń"
	reDzieCŃ = regexp.MustCompile("^(.+d)zie([cn])$")
	// Pasuje do rzeczowników zakończonych na -rzec,
	// np. "Marzec" i "Podgórzec"
	reRzec = regexp.MustCompile("^(.+r)ze(c)$")
	// Pasuje do takich rzeczowników zakończonych na -ek i -iek,
	// które mają co najmniej 2 sylaby, np. "budynek" i "Misiek"
	re2OrMoreSyllablesECK = regexp.MustCompile(
		"^(.*[aeiouy][^aeiouy])i?e([ck])$")
	// Pasuje do takich rzeczowników zakończonych na -el, -eł,
	// -eń, -er, -iel, -ieł, -ień i -ier, które mają co najmniej 2
	// sylaby, np. "Wróbel", "Styczeń", "magister", "Szczygieł",
	// "Wrzesień", "Węgier", "Nikiel", "inżynier"
	re2OrMoreSyllablesE = regexp.MustCompile(
		"^(.*[aeiouy][^aeiouy])i?e([lnr])$")
	// Pasuje do rzeczowników zakończonych na -ów lub -ach,
	// np. "Nowaków", "Kraków", "Nowakach" i "Stelmach"
	reÓwAch = regexp.MustCompile("^(.+)(ow|ach)$")
	// Pasuje do rzeczowników i przymiotników występujących w
	// takim przypadku i liczbie, w których końcówka zawiera
	// spółgłoskę, np. "Kowalskiego", "Mądrym", "profesorowi"
	reEndingsWithConsonant = regexp.MustCompile(
		"^(.+)(ego|emu|ej|im|ym|imi|ymi|ich|ych" +
			"|owi|owie|em|om|ami)$")
	// Pasuje do rzeczowników i przymiotników, których temat jest
	// dłuższy niż 1 znak, występujących w takim przypadku i
	// liczbie, w których końcówka składa się tylko z samogłosek,
	// np. "Kowalski", "geologii", "fizyka"
	reEndingsWithVowel = regexp.MustCompile("^(.+[^aeiouy])[aeiouy]*")
)

// ToASCIIStems usuwa z wyrazu `w` końcówkę rzeczownika lub
// przymiotnika
func ToASCIIStems(w ASCIIWord) []ASCIIStem {
	s := string(w)
	if m := reCieT.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1] + "t")
	}
	if m := reDzieD.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1] + "d")
	}
	if m := reCeK.FindStringSubmatch(s); m != nil {
		return nil // Zadanie 2
	}
	if m := reCyK.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1]+"c", m[1]+"k")
	}
	// Zadanie 3
	// Zadanie 4
	// Zadanie 5
	if m := reDzieCŃ.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(s, m[1]+m[2])
	}
	if m := reRzec.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1]+m[2], m[1]+"z"+m[2])
	}
	if m := re2OrMoreSyllablesECK.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1] + m[2])
	}
	if m := re2OrMoreSyllablesE.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(s, m[1]+m[2])
	}
	if m := reÓwAch.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(s, m[1])
	}
	if m := reEndingsWithConsonant.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(strings.TrimRight(m[1], "i"))
	}
	if m := reEndingsWithVowel.FindStringSubmatch(s); m != nil {
		return ToASCIIStemSlice(m[1])
	}
	return ToASCIIStemSlice(s)
}

// NonemptyStems usuwa puste łańcuchy z wycinka `stems`
func RemoveEmptyStems(stems []ASCIIStem) []ASCIIStem {
	ret := []ASCIIStem{}
	for _, s := range stems {
		if s != "" {
			ret = append(ret, s)
		}
	}
	return ret
}

// Częste spójniki, przyimki i skróty, które mogą mieć inne znaczenia,
// np. "I" jako numer piętra
var Stopwords = map[ASCIIStem]bool{
	"i":  true,
	"o":  true,
	"p":  true, // p. o. kierownika
	"w":  true,
	"z":  true,
	"ul": true,
}

// MeaningfulStems usuwa z wycinka tematów wyrazów `stems` takie
// łańcuchy, które można pomylić z innymi wyrazami lub z ich tematami
func RemoveStopwords(stems []ASCIIStem) []ASCIIStem {
	ret := []ASCIIStem{}
	for _, s := range stems {
		if !Stopwords[s] {
			ret = append(ret, s)
		}
	}
	return ret
}

// ASCIIStringToASCIIStemSlice dzieli łańcuch `as` na wyrazy i usuwa z
// tych wyrazów końcówki tak, że zostają tylko tematy wyrazów. Jeśli
// parametr `rmStopwords` ma wartość `true`, ASCIIStringToStemSlice
// usuwa takie tematy wyrazów, które można pomylić z innymi wyrazami
// lub z ich tematami
func ASCIIStringToASCIIStemSlice(as ASCIIString, rmStopwords bool) []ASCIIStem {
	ret := []ASCIIStem{}
	for _, s := range SplitASCIIString(as) {
		w := ToASCIIWord(s)
		stems := ToASCIIStems(w)
		stems = RemoveEmptyStems(stems)
		if rmStopwords {
			stems = RemoveStopwords(stems)
		}
		for _, stem := range stems {
			ret = append(ret, stem)
		}
	}
	return ret
}
