package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

type Client struct {
	client  http.Client
	Token   string
	Chat_id int
}

func New(c *config.Config) Client {
	return Client{
		Token:   c.TOKEN_ID,
		Chat_id: c.CHAT_ID,
		client:  http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) SendMessages(ctx context.Context, logger *slog.JSONHandler, message any) error {
	type content struct {
		Text    string `json:"text"`
		Chat_ID int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendMessage", config.URL, c.Token)

	d := content{
		Text:    message.(string),
		Chat_ID: c.Chat_id,
	}

	data, err := toMarshalAndBuffer(d)
	if data == nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to marshal body", err)
		return err
	}

	resp, err := http.Post(URL, "application/json", data)
	if err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to post message", err)
		return err
	}

	defer resp.Body.Close()

	var r map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to decode response body", err)
		return err
	}

	return nil
	//agregar fecha de la noticia al login para saber cuál fue la última
}

func (c *Client) SendPhoto(ctx context.Context, logger *slog.JSONHandler, post feed.Post) error {
	type content struct {
		Photo     string `json:"photo"`
		Caption   string `json:"caption"`
		ParseMode string `json:"parse_mode"`
		Chat_ID   int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendPhoto", config.URL, c.Token)

	caption := fmt.Sprintf("<a href='%v'>%v</a>", post.URL, post.Title)

	d := content{
		Photo:     post.PhotoURL,
		Chat_ID:   c.Chat_id,
		Caption:   caption,
		ParseMode: "HTML",
	}

	data, err := toMarshalAndBuffer(d)
	if data == nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to marshal body", err)
		return err
	}

	resp, err := http.Post(URL, "application/json", data)
	if err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to post message", err)
		return err
	}

	defer resp.Body.Close()

	var r map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to decode response body", err)
		return err
	}

	fmt.Println(r["ok"])

	return nil

}

func toMarshalAndBuffer(d any) (*bytes.Buffer, error) {
	p, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(p), nil
}
