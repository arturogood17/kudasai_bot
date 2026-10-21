package bot

func messageAuth(id, chat_id int) bool {
	return id == chat_id
}
