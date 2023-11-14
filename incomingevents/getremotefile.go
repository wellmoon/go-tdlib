package incomingevents

import "github.com/wellmoon/go-tdlib/entities"

type GetRemoteFile struct {
	Event

	*entities.File `json:"file"`
}
