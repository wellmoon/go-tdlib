package incomingevents

import "github.com/wellmoon/go-tdlib/entities"

type UpdateFile struct {
	Type string         `json:"@type"`
	File *entities.File `json:"file"`
}
