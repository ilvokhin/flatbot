package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	sent, err := readSent("sent/sent.json")
	if err != nil {
		log.Fatal(err)
	}
	flats = removeAlreadySent(flats, sent)
	m := messenger{
		Token:  os.Getenv("FLATBOT_TELEGRAM_BOT_API_TOKEN"),
		ChatID: os.Getenv("FLATBOT_TELEGRAM_CHANNEL_ID"),
	}
	for _, f := range flats {
		err = m.Send(f)
		if err != nil {
			// TODO: what to do with it?
			log.Print(err)
		}
		sent = append(sent, f)
	}
	fmt.Println(flats)
	writeSent(flats, "/tmp/sent.json")
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
