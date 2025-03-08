package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	filename := "2025-02-19-isle-of-dogs.html"
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Errorf("Could not read %v", filename)
	}
	want := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
		flat{
			ID:    158462822,
			Price: "£3,000",
		},
		flat{
			ID:    157948184,
			Price: "£2,400",
		}}
	got, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse failed: got: %v, want: %v", got, want)
	}
}
