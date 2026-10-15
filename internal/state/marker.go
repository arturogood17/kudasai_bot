// escribir y leer el marcador (líneas fáciles de identificar) y el caso de primera ejecución
package state

import (
	"log/slog"
	"time"
)

func Marker(logger *slog.Logger, title string, TFetch time.Time) {
	logger.Info("Marker Set", "Marker_date_set", time.Now(), "Last_title_fetch", title, "LastestFetch", TFetch)
}
