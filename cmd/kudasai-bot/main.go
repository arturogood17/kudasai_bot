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
		fmt.Printf("No se pudo abrir/crear el archivo de logging: %v", err)
		os.Exit(1)
	}

	ErrorFile, err := config.OpenOrCreateLogFile(config.ERRORFILEPATH)
	if err != nil {
		fmt.Printf("No se pudo abrir el archivo de logging de errores: %v", err)
		os.Exit(1)
	}

	defer LogFile.Close()
	defer ErrorFile.Close()

	logger := state.NewMultilogger(ErrorFile, LogFile)

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

	err = telegram.SendPosts(client, logger, filteredPosts)
	if err != nil {
		os.Exit(1)
	}

}
