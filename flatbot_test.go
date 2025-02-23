package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	filename := "htmls/2025-02-19-isle-of-dogs.html"
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Could not read %v", filename)
	}
	want := []flat{
		flat{
			URL:   "https://rightmove.co.uk/properties/156522206#",
			Price: "£2,500",
		},
		flat{
			URL:   "https://rightmove.co.uk/properties/158462822#",
			Price: "£3,000",
		},
		flat{
			URL:   "https://rightmove.co.uk/properties/157948184#",
			Price: "£2,400",
		}}
	got, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Parse failed: got: %v, want: %v", want, got)
	}
}
