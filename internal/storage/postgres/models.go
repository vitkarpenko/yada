package postgres

import (
	"gorm.io/gorm"
	"time"
)

type Reminder struct {
	gorm.Model
	MessageID string
	UserID    string
	ChannelID string
	RemindAt  time.Time
}
