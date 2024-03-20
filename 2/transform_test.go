package main

import (
	"slices"
	"testing"
)

func TestTransformPhoneNumbers(t *testing.T) {
	data := []struct {
		in   string
		want string
	}{
		{"Mój numer to +48 12 617 00 00", "Mój numer to 12-617-00-00"},
		{"126173456 i 126181235", "12-617-34-56 i 12-618-12-35"},
		{"+48  12  617  00  01", "12-617-00-01"},
		{"12 - 617 - 00 - 02", "12-617-00-02"},
		{"12-617-0003", "12-617-00-03"},
		{"12-6170004", "12-617-00-04"},
		{"910000005", "91-000-00-05"},
		{"Twój numer to +48 606 000 000", "Twój numer to 606-000-000"},
		{"+48  606  000  001", "606-000-001"},
		{"606 - 000 - 002", "606-000-002"},
		{"606 000-003", "606-000-003"},
		{"606 000004", "606-000-004"},
		{"606000005", "606-000-005"},
		{"Ich numer to +48 505 01 23 45", "Ich numer to 505-012-345"},
		{"+48  505  01  23  45", "505-012-345"},
		{"505 - 01 - 23 - 45", "505-012-345"},
		{"505 01-23-45", "505-012-345"},
	}
	for _, d := range data {
		if got := TransformPhoneNumbers(d.in); got != d.want {
			t.Errorf("NormalizePhoneNumber(%#v) == %#v want %#v",
				d.in, got, d.want)
		}
	}
}

func TestToASCIIString(t *testing.T) {
	in := "STRÓŻ pchnął KOŚĆ w QUIZ gędźb VEL fax MYJŃ."
	want := ASCIIString("stroz pchnal kosc w quiz gedzb vel fax myjn.")
	if got := ToASCIIString(in); got != want {
		t.Errorf("RemoveDiacritics(%#v) == %#v want %#v",
			in, got, want)
	}
}

func TestToASCIIWord(t *testing.T) {
	in := []ASCIIString{"Ala", ".:ma", "kota:.", ",,Elementarz''", ".", ""}
	want := []ASCIIWord{"Ala", "ma", "kota", "Elementarz", "", ""}
	for i := range in {
		if got := ToASCIIWord(in[i]); got != want[i] {
			t.Errorf("ToASCIIWord(%#v) == %#v want %#v",
				in[i], got, want[i])
		}
	}

}

func TestToASCIIStems(t *testing.T) {
	in := []ASCIIString{
		"nowak", "nowaka", "nowakowi", "nowakiem", "nowaku",
		"nowakowie", "nowakow",
		"nowakom", "nowakami", "nowakach",
		"wolski", "wolskiego", "wolskiemu", "wolskim",
		"wolscy", "wolskich", "wolskimi",
		"wolska", "wolskiej", "wolskie",
		"zimny", "zimnego", "zimnemu", "zimnym",
		"zimni", "zimnych", "zimnymi",
		"julia", "julii", "julie",
		"kot", "kocie", "alfred", "alfredzie",
		"agnieszka", "agnieszce", "kolega", "koledze",
		"doktor", "doktorze", "doktorzy",
		"jerzy", "jerzego",
		"marzec", "marca",
		"podgorzec", "podgorzca",
		"kwiecien", "kwietnia",
		"grudzien", "grudnia",
		"niemiec", "niemca", "bieniek", "bienka",
		"pawelec", "pawelca", "dudek", "dudka",
		"wrobel", "wrobla",
		"gerlach", "gerlachu",
		"dec", "deca", "piec", "pieca",
		"cwiek", "cwieka", "skrzek", "skrzeka",
	}
	want := [][]ASCIIStem{
		{"nowak"}, {"nowak"}, {"nowak"}, {"nowak"}, {"nowak"},
		{"nowak"}, {"nowakow", "nowak"},
		{"nowak"}, {"nowak"}, {"nowakach", "nowak"},
		{"wolsk"}, {"wolsk"}, {"wolsk"}, {"wolsk"},
		{"wolsc", "wolsk"}, {"wolsk"}, {"wolsk"},
		{"wolsk"}, {"wolsk"}, {"wolsk"},
		{"zimn"}, {"zimn"}, {"zimn"}, {"zimn"},
		{"zimn"}, {"zimn"}, {"zimn"},
		{"jul"}, {"jul"}, {"jul"},
		{"kot"}, {"kot"}, {"alfred"}, {"alfred"},
		{"agnieszk"}, {"agnieszk"}, {"koleg"}, {"koleg"},
		{"doktor"}, {"doktor", "doktorz"}, {"doktor", "doktorz"},
		{"jer", "jerz"}, {"jerz"},
		{"marc", "marzc"}, {"marc"},
		{"podgorc", "podgorzc"}, {"podgorzc"},
		{"kwiecien", "kwietn"}, {"kwietn"},
		{"grudzien", "grudn"}, {"grudn"},
		{"niemc"}, {"niemc"}, {"bienk"}, {"bienk"},
		{"pawelc"}, {"pawelc"}, {"dudk"}, {"dudk"},
		{"wrobel", "wrobl"}, {"wrobl"},
		{"gerlach", "gerl"}, {"gerlach"},
		{"dec"}, {"dec"}, {"piec"}, {"piec"},
		{"cwiek"}, {"cwiek"}, {"skrzek"}, {"skrzek"},
	}
	for i := range in {
		if got := ToASCIIStems(ASCIIWord(in[i])); !slices.Equal(got, want[i]) {
			t.Errorf("NormalizeWord(%#v) == %#v want %#v",
				in[i], got, want[i])
		}
	}
}

func TestRemoveEmptyStems(t *testing.T) {
	in := []ASCIIStem{"", "jan", "", "kowalsk"}
	want := []ASCIIStem{"jan", "kowalsk"}
	if got := RemoveEmptyStems(in); !slices.Equal(got, want) {
		t.Errorf("RemoveEmptyStems(%#v) == %#v want %#v",
			in, got, want)
	}
}

func TestRemoveStopwords(t *testing.T) {
	in := []ASCIIStem{"wydzial", "fizyk", "i", "informatyk", "stosowan"}
	want := []ASCIIStem{"wydzial", "fizyk", "informatyk", "stosowan"}
	if got := RemoveStopwords(in); !slices.Equal(got, want) {
		t.Errorf("RemoveStopwords(%v) == %v want %v", in, got, want)
	}
}

func TestASCIIStringToASCIIStemSlice(t *testing.T) {
	in := ASCIIString("w ,,elementarzu''  *ala*  ma kota i psa.")
	want := map[bool][]ASCIIStem{
		false: []ASCIIStem{
			"w", "elementarz", "al", "ma", "kot", "i", "ps"},
		true: []ASCIIStem{
			"elementarz", "al", "ma", "kot", "ps"},
	}
	for rmStopwords, w := range want {
		if got := ASCIIStringToASCIIStemSlice(in, rmStopwords); !slices.Equal(got, w) {
			t.Errorf("ASCIIStringToASCIIStemSlice(%v, %v) == %v want %v",
				in, rmStopwords, got, w)
		}
	}
}

/* Zadanie dodatkowe
func TestAbbreviateFacultyName(t *testing.T) {
	data := []struct{
		in string
		want [2]ASCIIStem
	}{
		{
			"Wydział Energetyki i Paliw, Biuro Administracyjne",
			[2]ASCIIStem{"weip", "wep"},
		},
		{
			"Wydział Odlewnictwa",
			[2]ASCIIStem{"wo", "wo"},
		},
		{
			"Wydział Fizyki i Informatyki Stosowanej, Dziekanat",
			[2]ASCIIStem{"wfiis", "wfis"},
		},
	}
	for _, d := range data {
		if got := [2]ASCIIStem{
			AbbreviateFacultyName(d.in, false),
			AbbreviateFacultyName(d.in, true),
		}; got != d.want {
			t.Errorf("AbbreviateFacultyName(%#v, false/true) == " +
				"%#v want %#v", d.in, got, d.want)
		}
	}
}
*/
