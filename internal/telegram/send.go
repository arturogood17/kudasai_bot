// envío de foto + título + link, y plan B sin imagen

package telegram

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func SendPosts(c Client, logger *slog.Logger, posts []feed.Post) error {
	numPosts := len(posts)
	i := 0
	for ; i <= numPosts-1; i++ {
		p := posts[i]
		if err := c.SendPhoto(logger, p); err != nil {
			logger.Warn(fmt.Sprintf("Failed to send (title) {%s}. Trying just text message", p.Title), "Warning_date", time.Now(), "Error", err)
			msg := fmt.Sprintf("<a href='%v'>%v</a>", p.URL, p.Title)
			if err = c.SendMessage(logger, msg); err != nil {
				state.NewLogError(logger, time.Now(), fmt.Sprintf("Sent %d/%d. Failed to send (title) {%s}", i, numPosts, p.Title), err)
			}
		}
		state.NewLogPostSuccess(logger, time.Now(), p.PubDate, fmt.Sprintf("Sent %d/%d. Last sent (title) {%s}", i+1, numPosts, p.Title)) //Debería ir en SendMessage y SendPhoto
	}

	state.Marker(logger, posts[i-1].Title, posts[i-1].PubDate) //Debería ir en SendMessage y SendPhoto

	return nil
}
