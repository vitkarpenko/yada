package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"yada/internal/storage/postgres"
)

const setRemindersCommand = "!remind"

var durations = map[string]time.Duration{
	"seconds": time.Second,
	"minutes": time.Minute,
	"hours":   time.Hour,
	"days":    24 * time.Hour,
	"weeks":   7 * 24 * time.Hour,
	"months":  30 * 24 * time.Hour,
	"years":   365 * 24 * time.Hour,
}

func (y *Yada) SetReminderHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if tokens[0] != setRemindersCommand {
		return
	}
	tokens = tokens[1:]

	if m.Message.MessageReference == nil {
		_, _ = y.Discord.ChannelMessageSend(
			m.ChannelID,
			"Ответь на сообщение которое я должен тебе припомнить.",
		)
		return
	}

	if len(tokens) != 2 {
		_, _ = y.Discord.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(
				"Некорректное количество аргументов. Делай, типа, так: `%s 2 days`",
				setRemindersCommand,
			),
		)
		return
	}

	countToken, durationToken := tokens[0], tokens[1]
	count, err := strconv.ParseInt(countToken, 10, 64)
	if err != nil || count <= 0 {
		_, _ = y.Discord.ChannelMessageSend(m.ChannelID, "Первый аргумент должен быть натуральным числом.")
		return
	}
	duration, ok := durations[durationToken]
	if !ok {
		_, _ = y.Discord.ChannelMessageSend(
			m.ChannelID,
			"Второй аргумент выбери из: seconds, minutes, hours, days, weeks, months, years.",
		)
		return
	}

	reminder := &postgres.Reminder{
		MessageID: m.MessageReference.MessageID,
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		RemindAt:  time.Now().Add(time.Duration(count * int64(duration))),
	}
	y.DB.Create(&reminder)
}

func (y *Yada) loadReminders() {
	y.DB.Find(&y.Reminders)
}

func (y *Yada) checkRemindersInBackground() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			y.checkReminders()
		}
	}()
}

func (y *Yada) checkReminders() {
	for _, reminder := range y.Reminders {
		if reminder.RemindAt.Before(time.Now()) {
			y.remind(reminder)
			y.DB.Delete(&reminder)
		}
	}
}

func (y *Yada) remind(reminder postgres.Reminder) {
	_, _ = y.Discord.ChannelMessageSendComplex(
		reminder.ChannelID,
		&discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s>, ты просил напомнить. :)", reminder.UserID),
			Reference: &discordgo.MessageReference{
				MessageID: reminder.MessageID,
				ChannelID: reminder.ChannelID,
				GuildID:   y.Config.GuildID,
			},
		},
	)
}
