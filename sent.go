package main

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
)

func readSent(filename string) ([]flat, error) {
	data, err := os.ReadFile(filename)
	sent := make([]flat, 0)
	if err != nil {
		// This is fine, as we might just started and didn't dump
		// anything into sent file yet.
		if errors.Is(err, os.ErrNotExist) {
			return sent, nil
		}
		return sent, err
	}
	err = json.Unmarshal(data, &sent)
	if err != nil {
		return sent, err
	}
	if !slices.IsSortedFunc(sent, compareID) {
		return nil, errors.New("Invalid sent: not sorted")
	}
	return sent, nil
}

func removeFlats(whenFound bool, from, superset []flat) []flat {
	if !slices.IsSortedFunc(superset, compareID) {
		slices.SortFunc(superset, compareID)
	}
	out := make([]flat, 0)
	for _, f := range from {
		_, found := slices.BinarySearchFunc(superset, f, compareID)
		if found == whenFound {
			continue
		}
		out = append(out, f)
	}
	return out
}

func removeAlreadySent(fetched []flat, sent []flat) []flat {
	whenFound := true
	return removeFlats(whenFound, fetched, sent)
}

func removeDelisted(sent []flat, allFlats []flat) []flat {
	whenFound := true
	return removeFlats(!whenFound, sent, allFlats)
}

func writeSent(sent []flat, filename string) error {
	if !slices.IsSortedFunc(sent, compareID) {
		slices.SortFunc(sent, compareID)
	}
	jsonData, err := json.Marshal(sent)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}
