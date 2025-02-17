package outgoingevents

import "github.com/wellmoon/go-tdlib/entities"

type LoadChats struct {
	ChatList *entities.ChatList `json:"chat_list"`
	Limit    int32              `json:"limit"`
}

func (s LoadChats) Type() string {
	return "loadChats"
}
