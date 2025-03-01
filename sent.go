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

func removeAlreadySent(fetched []flat, sent []flat) []flat {
	if !slices.IsSortedFunc(sent, compareID) {
		panic("Sent expected to be sorted")
	}
	recent := make([]flat, 0)
	for _, f := range fetched {
		_, found := slices.BinarySearchFunc(sent, f, compareID)
		if found {
			continue
		}
		recent = append(recent, f)
	}
	return recent
}

func writeSent(sent []flat, filename string) error {
	jsonData, err := json.Marshal(sent)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}
