package bot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func BotLogic(c telegram.Client, logger *slog.Logger, LogFile *os.File) error {
	URL := fmt.Sprintf("%v%v/getUpdates", config.URL, c.Token)

	type requestBody struct {
		Offset  int32 `json:"offset"`
		Timeout int32 `json:"timeout"`
	}

	client := http.Client{Timeout: 55 * time.Second}
	offset := 0

	for {

		reqB := requestBody{
			Offset:  int32(offset),
			Timeout: 45,
		}

		upd, err := requestAnddecoding(client, reqB, URL)
		if err != nil {
			state.NewLogError(logger, time.Now(), "Error trying to send the request", err)
			continue
		}

		if len(upd.Result) > 0 {
			mostRecentMsg := upd.Result[len(upd.Result)-1]
			offset = mostRecentMsg.UpdateID + 1
			if !messageAuth(mostRecentMsg.Message.From.ID, c.Chat_id) {
				c.SendMessages(logger, "Unavailable")
				state.NewLogError(logger, time.Now(), "Invalid user attempt", nil)
				continue
			}

			if strings.ToLower(mostRecentMsg.Message.Text) != "fetch" {
				c.SendMessages(logger, "Command not recongnize")
				continue
			}
		}
	}
}

func requestAnddecoding(c http.Client, reqB any, URL string) (Updates, error) {
	data, err := telegram.ToMarshalAndBuffer(reqB)
	if err != nil {
		return Updates{}, err
	}

	req, err := http.NewRequest("POST", URL, data)
	if err != nil {
		return Updates{}, err
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return Updates{}, err
	}

	defer resp.Body.Close()

	var upd Updates

	if err = json.NewDecoder(resp.Body).Decode(&upd); err != nil {
		return Updates{}, err
	}

	return upd, nil
}
