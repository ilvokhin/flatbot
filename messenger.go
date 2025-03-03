package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type messenger struct {
	Token  string
	ChatID string
}

func (m *messenger) Send(f flat) error {
	text := fmt.Sprintf("%v\n%v", f.Price, f.URL())
	data := url.Values{"chat_id": {m.ChatID}, "text": {text}}
	resp, err := http.PostForm(m.makeURL(), data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		type errorResponse struct {
			ErrorCode   int    `json:"error_code"`
			Description string `json:"description"`
		}
		errorResp := errorResponse{}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return err
		}
		return fmt.Errorf("Error code: %v, description: %v",
			errorResp.ErrorCode, errorResp.Description)
	}
	return nil
}

func (m *messenger) makeURL() string {
	return fmt.Sprintf(
		"https://api.telegram.org/bot%v/sendMessage",
		m.Token)
}
