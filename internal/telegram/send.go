// envío de foto + título + link, y plan B sin imagen

package telegram

import (
	"context"
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
)

func SendPosts(ctx context.Context, c Client, logger *slog.JSONHandler, posts []feed.Post) error {
	for _, p := range posts {
		if err := c.SendPhoto(ctx, logger, p); err != nil {
			return err
		}
	}

	return nil
}
