package bot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func BotLogic(c telegram.Client, logger *slog.Logger) {
	URL := fmt.Sprintf("%v%v/getUpdates", config.URL, c.Token)

	type requestBody struct {
		Offset  int32 `json:"offset"`
		Timeout int32 `json:"timeout"`
	}

	client := http.Client{Timeout: 55 * time.Second}
	offset := 0
	retryTime := 1 * time.Second
	countRetries := 0

	for {

		payload := requestBody{
			Offset:  int32(offset),
			Timeout: 45,
		}

		resp, err := telegram.RequestSending(&client, URL, "POST", payload)
		if err != nil {
			state.NewLogError(logger, time.Now(), "Error trying to send the request", err)
			time.Sleep(retryTime)
			retryTime += 1 * time.Second
			if retryTime >= 30*time.Second {
				retryTime = 30 * time.Second
			}
			continue
		}

		retryTime = 1 * time.Second

		var upd Updates

		if err = json.NewDecoder(resp.Body).Decode(&upd); err != nil {
			state.NewLogError(logger, time.Now(), "Error trying to decode the response from server", err)
			return
		}

		defer resp.Body.Close()

		if len(upd.Result) > 0 {
			mostRecentMsg := upd.Result[len(upd.Result)-1]
			offset = mostRecentMsg.UpdateID + 1
			if !messageAuth(mostRecentMsg.Message.From.ID, c.Chat_id) {
				c.SendMessage(logger, "Unavailable")
				state.NewLogError(logger, time.Now(), "Invalid user attempt", nil)
				continue
			}

			if strings.ToLower(mostRecentMsg.Message.Text) != "fetch" {
				c.SendMessage(logger, "Command not recongnize")
				continue
			}
			if err := CommandRun(logger, c); err != nil {
				state.NewLogError(logger, time.Now(), "Error serving posts", err)
				offset = mostRecentMsg.UpdateID
				time.Sleep(retryTime)
				retryTime += 1 * time.Second
				if retryTime >= 30*time.Second {
					retryTime = 30 * time.Second
				}
				countRetries += 1
				if countRetries >= 4 {
					offset = mostRecentMsg.UpdateID + 1
					countRetries = 0
				}
				continue
			}
			retryTime = 1 * time.Second
			countRetries = 0
		}
	}
}
