package entity

import "github.com/google/uuid"

type Student struct {
	ID               uuid.UUID
	FullName         string
	Phone            string
	TelegramChatID   int64
	TelegramUserName string
}
