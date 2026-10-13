// envío de foto + título + link, y plan B sin imagen

package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func SendPosts(ctx context.Context, c Client, logger *slog.JSONHandler, posts []feed.Post) error {
	numPosts := len(posts)
	i := 0
	for ; i <= numPosts-1; i++ {
		p := posts[i]
		if err := c.SendPhoto(ctx, logger, p); err != nil {
			msg := fmt.Sprintf("<a href='%v'>%v</a>", p.URL, p.Title)
			if err = c.SendMessages(ctx, logger, msg); err != nil {
				state.NewLogError(ctx, logger, time.Now(), slog.LevelError, fmt.Sprintf("Sent %d/%d. Failed to send (title) {%s}", i, numPosts, p.Title), err)
			}
		}
		state.NewLogPostSuccess(ctx, logger, time.Now(), p.PubDate, slog.LevelInfo, fmt.Sprintf("Sent %d/%d. Last sent (title) {%s}", i, numPosts, p.Title))
	}

	state.Marker(ctx, logger, posts[i-1].Title, posts[i-1].PubDate)

	return nil
}
