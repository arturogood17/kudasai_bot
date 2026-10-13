package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
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

	if configLoaded.CHAT_ID == 0 || configLoaded.TOKEN_ID == "" {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Failed to load one of the secrets", nil)
		os.Exit(1)
	}

	client := telegram.New(configLoaded)

	posts, err := feed.Feed(ctx, client.Clt, logger)
	if err != nil {
		os.Exit(1)
	}

	filteredPosts := filter.FilterByCategory(posts)

	filter.OrderPostsByDate(filteredPosts)

	err = telegram.SendPosts(ctx, client, logger, filteredPosts)
	if err != nil {
		os.Exit(1)
	}

}
