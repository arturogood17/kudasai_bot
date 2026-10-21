package bot

type Updates struct {
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			UpdateID  int `json:"update_id"`
			MessageID int `json:"message_id"`
			From      struct {
				ID int `json:"id"`
			} `json:"from"`
			Date int64  `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}
