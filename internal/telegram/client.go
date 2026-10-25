package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

type Client struct {
	Clt     *http.Client
	Token   string
	Chat_id int
}

func New(c *config.Config) Client {
	return Client{
		Token:   c.TOKEN_ID,
		Chat_id: c.CHAT_ID,
		Clt:     &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) SendMessage(logger *slog.Logger, message any) error {
	type content struct {
		Text    string `json:"text"`
		Chat_ID int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendMessage", config.URL, c.Token)

	payload := content{
		Text:    message.(string),
		Chat_ID: c.Chat_id,
	}

	resp, err := RequestSending(c.Clt, URL, "POST", payload)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error trying to send the request", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) SendPhoto(logger *slog.Logger, post feed.Post) error {
	type content struct {
		Photo     string `json:"photo"`
		Caption   string `json:"caption"`
		ParseMode string `json:"parse_mode"`
		Chat_ID   int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendPhoto", config.URL, c.Token)

	caption := fmt.Sprintf("<a href='%v'>%v</a>", post.URL, post.Title)

	if utf8.RuneCountInString(caption) > 1024 {
		caption = caption[:1024]
	}

	payload := content{
		Photo:     post.PhotoURL,
		Chat_ID:   c.Chat_id,
		Caption:   caption,
		ParseMode: "HTML",
	}

	resp, err := RequestSending(c.Clt, URL, "POST", payload)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error trying to send the request", err)
		return err
	}

	defer resp.Body.Close()

	return nil

}

func ToMarshalAndBuffer(d any) (*bytes.Buffer, error) {
	p, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(p), nil
}

func RequestSending(c *http.Client, URL string, method string, payload any) (*http.Response, error) {
	data, err := ToMarshalAndBuffer(payload)
	if data == nil {
		return nil, err
	}

	req, err := http.NewRequest(method, URL, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
