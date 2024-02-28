package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPodziel(t *testing.T) {
	data := []struct {
		in   string
		want [][]string
	}{
		{
			"dom",
			[][]string{
				{"d", "om"},
				{"do", "m"},
			},
		},
		{
			"kwiat",
			[][]string{
				{"k", "wiat"},
				{"kw", "iat"},
				{"kwi", "at"},
				{"kwia", "t"},
			},
		},
	}
	for _, d := range data {
		got := Split(d.in)
		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", d.want) {
			t.Errorf("got %#v; want %#v", got, d.want)
		}
	}
}

func TestAddIfIn(t *testing.T) {
	words := map[string]bool{
		"akt": true,
		"or":  true,
	}
	data := []struct {
		in   []string
		want map[string]int
	}{
		{
			[]string{"los", "owo"},
			map[string]int{},
		},
		{
			[]string{"akt", "or"},
			map[string]int{"akt": 1, "or": 1},
		},
	}
	for _, d := range data {
		got := map[string]int{}
		AddIfIn(d.in, words, &got)
		if !reflect.DeepEqual(got, d.want) {
			t.Errorf("got %#v; want %#v", got, d.want)
		}
	}
}

func TestSort(t *testing.T) {
	data := []struct {
		in   map[string]int
		want []Pair
	}{
		{
			map[string]int{},
			[]Pair{},
		},
		{
			map[string]int{"opowiem": 2, "uparty": 3, "określ": 2},
			[]Pair{Pair{"uparty", 3}, {"określ", 2}, {"opowiem", 2}},
		},
	}
	for _, d := range data {
		got := Sort(d.in)
		if !reflect.DeepEqual(got, d.want) {
			t.Errorf("got %#v; want %#v", got, d.want)
		}
	}
}
