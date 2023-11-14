package incomingevents

import "github.com/wellmoon/go-tdlib/entities"

type UpdateNewMessage struct {
	Type    string            `json:"@type"`
	Message *entities.Message `json:"message"`
}
