package outgoingevents

import "github.com/wellmoon/go-tdlib/entities"

type GetCallbackQueryAnswer struct {
	ChatID    int64                          `json:"chat_id"`
	MessageID int64                          `json:"message_id"`
	Payload   *entities.CallbackQueryPayload `json:"payload"`
}

func (s GetCallbackQueryAnswer) Type() string {
	return "getCallbackQueryAnswer"
}
