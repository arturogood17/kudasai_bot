// se queda solo con Anime y Manga
package filter

import (
	"github.com/arturogood17/kudasai_bot/internal/feed"
)

const Manga = "Manga"
const Anime = "Anime"

func FilterByCategory(posts []feed.Post) []feed.Post {
	var filteredPosts []feed.Post
	for _, p := range posts {
		if p.Category == Manga || p.Category == Anime {
			filteredPosts = append(filteredPosts, p)
		}
	}
	return filteredPosts
}
