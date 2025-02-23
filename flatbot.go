package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := "http://localhost:8000/2025-02-19-isle-of-dogs.html"
	body, err := fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	flats, err := parse(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(flats)
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), err
	}
	if resp.StatusCode != http.StatusOK {
		return make([]byte, 0),
			fmt.Errorf("Bad response status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), err
	}
	return body, nil
}

type flat struct {
	URL   string
	Price string
}

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
		attr := matchAttr(n, "data-testid")
		if attr == nil || attr.Val != "property-price" {
			continue
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
	f := flat{URL: makeURL(url.Val), Price: ""}
	for n := range root.Descendants() {
		if price, found := strings.CutSuffix(n.Data, " pcm"); found {
			f.Price = price
			return f, nil
		}
	}
	return flat{}, errors.New("Couldn't find price")
}

func makeURL(path string) string {
	prettySuffix, _ := strings.CutSuffix(path, "/?channel=RES_LET")
	return fmt.Sprintf("https://rightmove.co.uk%v", prettySuffix)
}
