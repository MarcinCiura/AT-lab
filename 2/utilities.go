package main

import (
	"fmt"
	"strings"
	"unsafe"
)

type (
	// Łańcuch małych liter bez znaków diakrytycznych, na przykład
	// "ola ma zolwia"
	ASCIIString string
	// Wyraz zapisany małymi lterami bez znaków diakrytycznych, na
	// przykład "zolwia"
	ASCIIWord string
	// Temat wyrazu zapisany małymi literami bez znaków
	// diakrytycznych, na przykład "zolw"
	ASCIIStem string
)

// ToStringSlice zmienia typ wycinka `sl` na wycinek łańcuchów
func ToStringSlice[S interface{ ~string }](sl []S) []string {
	return *(*[]string)(unsafe.Pointer(&sl))
}

// ToASCIIStemSlice tworzy z łańcuchów `ss` wycinek wartości typu
// `ASCIIStem`
func ToASCIIStemSlice(ss ...string) []ASCIIStem {
	return *(*[]ASCIIStem)(unsafe.Pointer(&ss))
}

// ToAnySlice tworzy z wycinka `sl` wycinek wartości typu `any`
func ToAnySlice(sl []ColumnName) []any {
	ret := []any{}
	for _, c := range sl {
		ret = append(ret, c)
	}
	return ret
}

// MakeStringSliceAndAnySlice tworzy wycinek `n` łańcuchów i wycinek
// `n` wartości typu `any`. Każdy element tego drugiego wycinka
// wskazuje na odpowiedni element pierwszego wycinka
func MakeStringSliceAndAnySlice(n int) ([]string, []any) {
	ss := make([]string, n)
	as := make([]any, n)
	for i := range ss {
		as[i] = &ss[i]
	}
	return ss, as
}

// SplitASCIIString dzieli łańcuch `s` na takie części, pomiędzy
// którymi leżą dowolne ciągi 1 lub więcej białych znaków
func SplitASCIIString(s ASCIIString) []ASCIIString {
	ret := []ASCIIString{}
	for _, p := range strings.Fields(string(s)) {
		ret = append(ret, ASCIIString(p))
	}
	return ret
}

// JoinASCIIStems łączy elementy wycinka `s` spacjami
func JoinASCIIStems(s []ASCIIStem) string {
	return strings.Join(ToStringSlice(s), " ")
}

// JoinQuotedStems otacza elementy wycinka `ss` cudzysłowami, łączy te
// elementy kopiami łańcucha `joiner` i otacza wynik nawiasami
// okrągłymi.
//
// Przykład:
// JoinQuotedStems(
//
//	[]ASCIIStem{"kij", "kijow"}, " OR ") == `("kij" OR "kijow")`
func JoinQuotedStems(ss []ASCIIStem, joiner string) string {
	ret := []string{}
	for _, s := range ss {
		ret = append(ret, fmt.Sprintf(`"%s"`, s))
	}
	return fmt.Sprintf(`(%s)`, strings.Join(ret, joiner))
}
