package main

import (
	"fmt"
	"os"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func main() {

	LogFile, err := config.OpenOrCreateLogFile(config.LOGFILEPATH)
	if err != nil {
		fmt.Printf("Couldn't open/create the logging file: %v", err)
		os.Exit(1)
	}

	ErrorFile, err := config.OpenOrCreateLogFile(config.ERRORFILEPATH)
	if err != nil {
		fmt.Printf("Couldn't open/create the error logging file: %v", err)
		os.Exit(1)
	}

	defer LogFile.Close()
	defer ErrorFile.Close()

	logger := state.NewMultilogger(ErrorFile, LogFile)

	markerNews, err := filter.GetMarker(logger, LogFile)
	if err != nil {
		os.Exit(1)
	}

	configLoaded := config.Load(logger)
	if configLoaded == nil {
		os.Exit(1)
	}

	if configLoaded.CHAT_ID == 0 || configLoaded.TOKEN_ID == "" {
		state.NewLogError(logger, time.Now(), "Failed to load one of the secrets", nil)
		os.Exit(1)
	}

	client := telegram.New(configLoaded)

	posts, err := feed.Feed(client.Clt, logger)
	if err != nil {
		os.Exit(1)
	}

	filteredPosts := filter.FilterByCategory(posts)

	filter.OrderPostsByDate(filteredPosts)

	lastestPosts := filter.FilterByDate(logger, markerNews, filteredPosts)
	if len(lastestPosts) == 0 {
		logger.Info("Nothing to fetch")
		os.Exit(0)
	}

	err = telegram.SendPosts(client, logger, lastestPosts)
	if err != nil {
		os.Exit(1)
	}

}
