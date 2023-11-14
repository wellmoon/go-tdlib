package incomingevents

import "github.com/wellmoon/go-tdlib/entities"

type GetMe struct {
	Event

	*entities.User `json:"user"`
}
