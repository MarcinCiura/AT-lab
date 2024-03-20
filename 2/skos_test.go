package main

import (
	"reflect"
	"testing"
)

func TestParseQuestion(t *testing.T) {
	data := []struct {
		question ASCIIString
		match    string
		cols     []ColumnName
	}{
		{
			"jacy profesorowie pracuja w budynku c-1?",
			`("prof") AND ("c-1")`,
			[]ColumnName{Osoba, Budynek},
		},
		{
			"kto ma telefon o numerze 12-617-12-34?",
			`("12-617-12-34")`,
			[]ColumnName{Osoba, Telefon},
		},
		{
			"czy znasz kogos o imieniu filip?",
			`("filip")`,
			[]ColumnName{Osoba},
		},
		{
			"którzy doktorzy z wydziału chemii nie sa habilitowani?",
			`("dr") AND ("chem") NOT ("hab")`,
			[]ColumnName{Osoba, Jednostka},
		},
		{
			"jakich znasz magistrow nowakow w budynku c-1?",
			`("mgr") AND ("nowakow" OR "nowak") AND ("c-1")`,
			[]ColumnName{Osoba, Budynek},
		},
		{
			"kto nie jest habilitowany?",
			"",
			nil,
		},
	}
	dbStems := map[ASCIIStem]map[ColumnName]bool{
		"mgr":          map[ColumnName]bool{Osoba: true},
		"dr":           map[ColumnName]bool{Osoba: true},
		"hab":          map[ColumnName]bool{Osoba: true},
		"prof":         map[ColumnName]bool{Osoba: true},
		"filip":        map[ColumnName]bool{Osoba: true},
		"nowak":        map[ColumnName]bool{Osoba: true},
		"nowakow":      map[ColumnName]bool{Osoba: true},
		"wydzial":      map[ColumnName]bool{Jednostka: true},
		"chem":         map[ColumnName]bool{Jednostka: true},
		"c-1":          map[ColumnName]bool{Budynek: true},
		"12-617-12-34": map[ColumnName]bool{Telefon: true},
	}
	for _, d := range data {
		match, cols, _ := ParseQuestion(d.question, dbStems)
		if match != d.match || !reflect.DeepEqual(cols, d.cols) {
			t.Errorf("ParseQuestion(%#v) == %#v, %#v want %#v, %#v",
				d.question, match, cols, d.match, d.cols)
		}
	}
}

func TestMakeQuery(t *testing.T) {
	match := `("mgr") AND ("nowakow" OR "nowak") AND ("c-1")`
	cols := []ColumnName{Osoba, Budynek}
	want := `SELECT osoba, budynek
FROM Pracownicy JOIN PracownicyFTS
ON Pracownicy.rowid = PracownicyFTS.rowid
WHERE dane MATCH '("mgr") AND ("nowakow" OR "nowak") AND ("c-1")'`
	if got := MakeQuery(match, cols); got != want {
		t.Errorf("MakeQuery(%#v, %#v) == %#v want %#v",
			match, cols, got, want)
	}
}
