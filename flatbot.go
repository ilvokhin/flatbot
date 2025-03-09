package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	apiToken = flag.String("api-token",
		os.Getenv("FLATBOT_TELEGRAM_BOT_API_TOKEN"),
		"Telegram Bot API token")
	chatID = flag.String("chat-id",
		os.Getenv("FLATBOT_TELEGRAM_CHAT_ID"),
		"Telegram chat id where to send notification messages")
	state = flag.String("state", "/tmp/flatbot-sent.json",
		"Filename to save and load already sent flats")
	interval = flag.Duration("frequency", 5*time.Minute,
		"Frequency interval to fetch new data")
	once   = flag.Bool("once", false, "Run fetch and message loop only once")
	dryRun = flag.Bool("dry-run", false,
		"Run entire flow, but print new flats to stdout instead of "+
			"sending them to Telegram")
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: flatbot [-api-token token] [-chat-id id] "+
			"[-state file] [-interval duration] [-once] [-send] "+
			"URL...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if !*dryRun {
		if len(*apiToken) == 0 {
			fmt.Fprintf(os.Stderr,
				"Going to send messages to Telegram, "+
					"but no API token was provided\n")
			os.Exit(1)
		}
		if len(*chatID) == 0 {
			fmt.Fprintf(os.Stderr,
				"Going to send messages to Telegram, "+
					"but no chat id was provided\n")
			os.Exit(1)
		}
	}
	if flag.NArg() == 0 {
		usage()
		os.Exit(1)
	}

	for {
		err := loopOnce()
		if err != nil {
			log.Fatal(err)
		}
		if *once {
			break
		}
		log.Printf("Going to sleep for %v", *interval)
		time.Sleep(*interval)
	}
}

func loopOnce() error {
	sent, err := readSent(*state)
	if err != nil {
		return err
	}
	fetched := make([]flat, 0)
	for _, url := range flag.Args() {
		sent, fetched, err = doOneURL(url, sent, fetched)
		if err != nil {
			log.Print(err)
			continue
		}
	}
	sent = removeDelisted(sent, fetched)
	return writeSent(sent, *state)
}

func doOneURL(url string, sent, fetched []flat) ([]flat, []flat, error) {
	body, err := fetch(url)
	if err != nil {
		return sent, fetched, err
	}
	parsed, err := parse(body)
	if err != nil {
		return sent, fetched, err
	}
	fetched = append(fetched, parsed...)
	newFlats := removeAlreadySent(parsed, sent)
	m := messenger{
		Token:  *apiToken,
		ChatID: *chatID,
	}
	for _, f := range newFlats {
		if !*dryRun {
			err = m.Send(f)
			if err != nil {
				log.Print(err)
				continue
			}
		}
		log.Printf("Should have been sent to chat: %v, %v",
			f.Price, f.URL())
		sent = append(sent, f)
	}
	return sent, fetched, nil
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
