// lectura del XML y conversión de fechas a UTC
package feed

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/state"
)

func PubDateParser(ctx context.Context, logger *slog.JSONHandler, pubDate string, title string) time.Time {
	t, err := time.Parse(time.RFC1123Z, pubDate)
	if err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, fmt.Sprintf("Error parsing date for %s", title), err)
		return time.Time{}
	}
	return t
}
