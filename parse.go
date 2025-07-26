package main

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func parse(body []byte) ([]flat, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return make([]flat, 0), err
	}
	flats := make([]flat, 0)
	for _, n := range findNodes(doc) {
		flat, err := parseNode(n)
		if err != nil {
			continue
		}
		flats = append(flats, flat)
	}
	slices.SortFunc(flats, compareID)
	flats = slices.CompactFunc(flats, func(a, b flat) bool {
		return compareID(a, b) == 0
	})
	return flats, nil
}

func findNodes(root *html.Node) []*html.Node {
	flats := make([]*html.Node, 0)
	for n := range root.Descendants() {
		if n.Type != html.ElementNode {
			continue
		}
		if n.Data != "a" {
			continue
		}
		// Try to match attribute for let.
		attr := matchAttr(n, "data-testid")
		if attr == nil || attr.Val != "property-price" {
			// If unsuccessful, try to match attribute for buy.
			attr = matchAttr(n, "data-test")
			if attr == nil || attr.Val != "property-header" {
				continue
			}
			href := matchAttr(n, "href")
			if href == nil || href.Val == "" {
				continue
			}
		}
		flats = append(flats, n)
	}
	return flats
}

func matchAttr(n *html.Node, key string) *html.Attribute {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return &attr
		}
	}
	return nil
}

func parseNode(root *html.Node) (flat, error) {
	url := matchAttr(root, "href")
	if url == nil {
		return flat{}, errors.New("Couldn't find URL")
	}
	ID, err := parseID(url.Val)
	if err != nil {
		return flat{}, err
	}
	f := flat{ID, ""}
	for n := range root.Descendants() {
		if price, found := strings.CutSuffix(n.Data, " pcm"); found {
			f.Price = price
			return f, nil
		}
		if strings.HasPrefix(n.Data, "£") {
			f.Price = strings.TrimSpace(n.Data)
			return f, nil
		}
	}
	return flat{}, errors.New("Couldn't find price")
}

func parseID(path string) (int, error) {
	s, _ := strings.CutPrefix(path, "/properties/")
	for _, channel := range []string{"LET", "BUY"} {
		suffix := fmt.Sprintf("#/?channel=RES_%v", channel)
		maybeID, _ := strings.CutSuffix(s, suffix)
		ID, err := strconv.Atoi(maybeID)
		if err == nil {
			return ID, nil
		}
	}
	err := fmt.Errorf("Couldn't extract ID from %q", path)
	return -1, err
}
