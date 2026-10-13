// escribir y leer el marcador (líneas fáciles de identificar) y el caso de primera ejecución
package state

import (
	"context"
	"log/slog"
	"time"
)

func Marker(ctx context.Context, logger *slog.JSONHandler, title string, TFetch time.Time) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "Marker Set", 0)
	record.Add("Title", title)
	record.Add("LastestFetch", TFetch)
	logger.Handle(ctx, record)
}
