package config

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/state"
	viper "github.com/spf13/viper"
)

const URL = "https://api.telegram.org/bot"
const FEED = "https://somoskudasai.com/noticias/feed/"
const ERRORFILEPATH = "ErrorFile.log"
const LOGFILEPATH = "LogInfoFile.log"

type Config struct {
	TOKEN_ID string
	CHAT_ID  int
}

func Load(logger *slog.Logger) *Config {
	var config *Config
	viper := viper.New()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			state.NewLogError(logger, time.Now(), "Unable to read secrets", err)
			return nil
		}
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to unmarshal", err)
		return nil
	}
	return config
}

func OpenOrCreateLogFile(filename string) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	FileCreation := path.Join(cwd, "data", filename)
	File, err := os.OpenFile(FileCreation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return File, nil
}
