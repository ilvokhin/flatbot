package main

import (
	"cmp"
	"fmt"
)

type flat struct {
	ID    int
	Price string
}

func (f *flat) URL() string {
	return fmt.Sprintf("https://rightmove.co.uk/properties/%v", f.ID)
}

func compareID(a, b flat) int {
	return cmp.Compare(a.ID, b.ID)
}
