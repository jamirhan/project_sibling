package tclient

type ChatID int64

type User struct {
	ID int64 `json:"id"`
}

type Chat struct {
	ID ChatID `json:"id"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Message struct {
	From            User            `json:"from"`
	Chat            Chat            `json:"chat"`
	Text            string          `json:"text"`
	MessageEntities []MessageEntity `json:"entities"`
	ReplyToMessage  *Message        `json:"reply_to_message"`
	MessageID int64 `json:"message_id"`
}

type ReplyParams struct {
	MessageID int64 `json:"message_id"`
}

type sendMessageRequest struct {
	ChatID ChatID `json:"chat_id"`
	Text   string `json:"text"`
	ReplyParameters ReplyParams `json:"reply_parameters"`
}
