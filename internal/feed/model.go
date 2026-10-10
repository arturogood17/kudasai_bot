// estructuras XML (título, link, categorías, fecha, media:content)

package feed

import "time"

type Post struct {
	Date      time.Time
	Title     string
	URL       string
	PhotoURL  string
	Category  string
	MediaType string
}
