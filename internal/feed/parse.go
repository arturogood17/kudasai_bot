// lectura del XML y conversión de fechas a UTC
package feed

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/state"
)

func PubDateParser(logger *slog.Logger, pubDate string, title string) time.Time {
	t, err := time.Parse(time.RFC1123Z, pubDate)
	if err != nil {
		state.NewLogError(logger, time.Now(), fmt.Sprintf("Error parsing date for %s", title), err)
		return time.Time{}
	}
	return t
}
