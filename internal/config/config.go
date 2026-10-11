package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/state"
	viper "github.com/spf13/viper"
)

const URL = "https://api.telegram.org/bot"
const FEED = "https://somoskudasai.com/noticias/feed/"

type Config struct {
	TOKEN_ID string
	CHAT_ID  int
}

func Load(ctx context.Context, logger *slog.JSONHandler) *Config {
	var config *Config
	viper := viper.New()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to read secrets", err)
			return nil
		}
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		state.NewLogError(ctx, logger, time.Now(), slog.LevelError, "Unable to unmarshal", err)
		return nil
	}
	return config
}
