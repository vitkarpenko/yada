package postgres

import (
	"time"
)

type Reminder struct {
	ID        uint `gorm:"primarykey"`
	MessageID string
	UserID    string
	ChannelID string
	RemindAt  time.Time
}

type LastQuote struct {
	ID   uint `gorm:"primarykey"`
	Hash string
}
