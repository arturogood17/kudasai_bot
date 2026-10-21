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

func (c *Client) SendMessages(logger *slog.Logger, message any) error {
	type content struct {
		Text    string `json:"text"`
		Chat_ID int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendMessage", config.URL, c.Token)

	d := content{
		Text:    message.(string),
		Chat_ID: c.Chat_id,
	}

	data, err := ToMarshalAndBuffer(d)
	if data == nil {
		state.NewLogError(logger, time.Now(), "Unable to marshal body", err)
		return err
	}

	resp, err := http.Post(URL, "application/json", data)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to post message", err)
		return err
	}

	defer resp.Body.Close()

	//For debugging. Needs to be erased.
	var r map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		state.NewLogError(logger, time.Now(), "Unable to decode response body", err)
		return err
	}

	return nil
	//agregar fecha de la noticia al login para saber cuál fue la última
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

	d := content{
		Photo:     post.PhotoURL,
		Chat_ID:   c.Chat_id,
		Caption:   caption,
		ParseMode: "HTML",
	}

	data, err := ToMarshalAndBuffer(d)
	if data == nil {
		state.NewLogError(logger, time.Now(), "Unable to marshal body", err)
		return err
	}

	resp, err := http.Post(URL, "application/json", data)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to post message", err)
		return err
	}

	defer resp.Body.Close()

	//For debugging. Needs to be erased.

	var r map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		state.NewLogError(logger, time.Now(), "Unable to decode response body", err)
		return err
	}

	return nil

}

func ToMarshalAndBuffer(d any) (*bytes.Buffer, error) {
	p, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(p), nil
}
