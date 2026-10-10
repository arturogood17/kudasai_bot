package main

import (
	"context"
	"os"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func main() {

	ctx := context.Background()

	logger := state.NewJSONLogger()

	configLoaded := config.Load(ctx, logger)
	if configLoaded == nil {
		os.Exit(1)
	}

	client := telegram.New(configLoaded)

	if err := client.SendPhoto(ctx, logger, feed.Post{
		Title:    "El manhwa The Remarried Empress confirma su adaptación al anime para 2027",
		URL:      "https://somoskudasai.com/noticias/el-manhwa-the-remarried-empress-confirma-su-adaptacion-al-anime-para-2027",
		PhotoURL: "https://cdn.somoskudasai.com/width=1280,height=720,quality=75,format=auto,fit=cover/2026/10/the-remarried-empress-anime.jpg",
		Category: "Anime",
	}); err != nil {
		os.Exit(1)
	}
}
