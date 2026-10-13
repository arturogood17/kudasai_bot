// petición HTTP al feed con timeout y comprobación del código 200
package feed

import (
	"context"
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func Feed(ctx context.Context, c *http.Client, logger *slog.JSONHandler) ([]Post, error) {
	var rss RSS
	logger.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "Feetching started", 0))
	resp, err := c.Get(config.FEED)
	if err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to obtain Kudasai feed", err)
		return nil, err
	}
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable parse xml file", err)
		return nil, err
	}

	defer resp.Body.Close()

	return Posts(ctx, logger, rss), err

}

func Posts(ctx context.Context, logger *slog.JSONHandler, r RSS) []Post {
	var P []Post
	for _, post := range r.Channel.Item {
		np := Post{
			Title:    post.Title,
			URL:      post.Link,
			PubDate:  PubDateParser(ctx, logger, post.PubDate, post.Title),
			Category: post.Category,
			PhotoURL: post.Content.URL,
		}
		P = append(P, np)
	}
	return P
}
