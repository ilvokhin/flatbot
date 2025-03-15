package main

import (
	"cmp"
	"fmt"
)

type flat struct {
	ID    int    `json:"id"`
	Price string `json:"price"`
}

func (f *flat) URL() string {
	return fmt.Sprintf("https://rightmove.co.uk/properties/%v", f.ID)
}

func compareID(a, b flat) int {
	return cmp.Compare(a.ID, b.ID)
}
