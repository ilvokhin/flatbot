package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseBasic(t *testing.T) {
	filename := "2025-02-19-basic.html"
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("Could not read %v", filename)
	}
	want := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
		flat{
			ID:    157948184,
			Price: "£2,400",
		},
		flat{
			ID:    158462822,
			Price: "£3,000",
		}}
	got, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse failed: got: %v, want: %v", got, want)
	}
}

func TestParseDulicates(t *testing.T) {
	filename := "2025-03-17-duplicates.html"
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("Could not read %v", filename)
	}
	want := []flat{
		flat{
			ID:    158595710,
			Price: "£2,000",
		},
		flat{
			ID:    158825903,
			Price: "£2,500",
		},
		flat{
			ID:    159476474,
			Price: "£3,000",
		},
		flat{
			ID:    159479504,
			Price: "£890",
		}}
	got, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse failed: got: %v, want: %v", got, want)
	}
}
