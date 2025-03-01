package main

import (
	"testing"
)

func TestURL(t *testing.T) {
	f := flat{ID: 0, Price: "£1,000"}
	got := f.URL()
	want := "https://rightmove.co.uk/properties/0"
	if got != want {
		t.Errorf("URL call failed: got: %v, want: %v", got, want)
	}
}

func TestCompareID(t *testing.T) {
	a := flat{ID: 0, Price: "£3,000"}
	b := flat{ID: 1, Price: "£1,000"}
	got := compareID(a, b)
	if got != -1 {
		t.Errorf("Wrong compare result: got: %v, want: %v", got, -1)
	}
	got = compareID(b, a)
	if got != 1 {
		t.Errorf("Wrong compare result: got: %v, want: %v", got, 1)
	}
	got = compareID(a, a)
	if got != 0 {
		t.Errorf("Wrong compare result: got: %v, want: %v", got, 0)
	}
}
