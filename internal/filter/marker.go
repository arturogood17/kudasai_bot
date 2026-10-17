// descarta las noticias anteriores o iguales al marcador y ordena de más antigua a más reciente
package filter

import (
	"encoding/json"
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func OrderPostsByDate(posts []feed.Post) {
	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].PubDate.Before(posts[j].PubDate)
	})
}

func readLogs(jsonFile *os.File) ([]map[string]any, error) {
	var fullFile []map[string]any
	decoder := json.NewDecoder(jsonFile)

	for decoder.More() {
		var data map[string]any
		err := decoder.Decode(&data)
		if err != nil {
			return nil, err
		}
		fullFile = append(fullFile, data)
	}
	return fullFile, nil
}

func GetMarker(logger *slog.Logger, logs *os.File) (map[string]any, error) {
	LogNews, err := readLogs(logs)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error obtaining the marker", err)
		return nil, err
	}

	var m map[string]any

	if len(LogNews) != 0 {
		for i, v := range LogNews {
			if _, ok := LogNews[i]["Last_title_fetch"]; ok {
				if _, ok = LogNews[i]["LastestFetch"]; ok {
					m = v
				}
			}
		}
	} else {
		return nil, nil
	}
	return m, nil
}

func FilterByDate(logger *slog.Logger, m map[string]any, posts []feed.Post) []feed.Post {
	if len(m) == 0 {
		return posts
	}
	lastestFetch := m["LastestFetch"]
	t := feed.PubDateParser(logger, lastestFetch.(string), "Filtering By Date", time.RFC3339)
	var lastestPost []feed.Post
	for _, p := range posts {
		if t.Before(p.PubDate) {
			lastestPost = append(lastestPost, p)
		}
	}
	return lastestPost
}
