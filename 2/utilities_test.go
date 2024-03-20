package main

import (
	"slices"
	"testing"
)

func TestToStringSlice(t *testing.T) {
	in := []ASCIIStem{"al", "m", "kot"}
	want := []string{"al", "m", "kot"}
	if got := ToStringSlice(in); !slices.Equal(got, want) {
		t.Errorf("ToStringSlice(%#v) == %#v want %#v", in, got, want)
	}
}

func TestToASCIIStemSlice(t *testing.T) {
	want := []ASCIIStem{"al", "m", "kot"}
	if got := ToASCIIStemSlice("al", "m", "kot"); !slices.Equal(got, want) {
		t.Errorf(`ToASCIIStemSlice("al", "m", "kot") == %#v want %#v`,
			got, want)
	}
}

func TestToAnySlice(t *testing.T) {
	in := []ColumnName{Osoba, Jednostka}
	want := []any{Osoba, Jednostka}
	if got := ToAnySlice(in); !slices.Equal(got, want) {
		t.Errorf(`ToAnySlice(%#v) == %#v want %#v`, in, got, want)
	}
}

func TestMakeStringSliceAndAnySlice(t *testing.T) {
	ss, as := MakeStringSliceAndAnySlice(2)
	if len(ss) != 2 || len(as) != 2 {
		t.Errorf("MakeStringSliceAndAnySlice(2) == %#v, %#v want len(slices) == 2",
			ss, as)
	}
	if as[0] != &ss[0] || as[1] != &ss[1] {
		t.Errorf("MakeStringSliceAndAnySlice(2) == []string{*%#v *%#v}, %#v",
			&ss[0], &ss[1], as)
	}
}

func TestSplitASCIIString(t *testing.T) {
	in := ASCIIString("  Ala  ma\tkota")
	want := []ASCIIString{"Ala", "ma", "kota"}
	if got := SplitASCIIString(in); !slices.Equal(got, want) {
		t.Errorf("SplitASCIIString(%#v) == %#v want %#v",
			in, got, want)
	}
}

func TestJoinASCIIStems(t *testing.T) {
	in := []ASCIIStem{"al", "m", "kot"}
	want := "al m kot"
	if got := JoinASCIIStems(in); got != want {
		t.Errorf("JoinASCIIStems(%v) == %#v want %#v", in, got, want)
	}
}

func TestJoinQuotedStems(t *testing.T) {
	in := []ASCIIStem{"marc", "marzc"}
	want := `("marc" OR "marzc")`
	if got := JoinQuotedStems(in, ` OR `); got != want {
		t.Errorf("JoinQuotedStems(%v) == %v want %v", in, got, want)
	}
}
