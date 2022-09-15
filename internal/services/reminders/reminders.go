package reminders

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vitkarpenko/yada/internal/config"
	"github.com/vitkarpenko/yada/storages/sqlite"
)

const (
	setRemindersCommand = "!remind"
	checkInterval       = time.Second
)

var durations = map[string]time.Duration{
	"seconds": time.Second,
	"minutes": time.Minute,
	"hours":   time.Hour,
	"days":    24 * time.Hour,
	"weeks":   7 * 24 * time.Hour,
	"months":  30 * 24 * time.Hour,
	"years":   365 * 24 * time.Hour,
}

type Service struct {
	discord *discordgo.Session
	queries *sqlite.Queries

	config config.Config

	reminders []sqlite.Reminder
}

func New(discord *discordgo.Session, queries *sqlite.Queries, config config.Config) *Service {
	service := &Service{discord: discord, queries: queries, config: config}
	service.loadReminders()
	service.checkInBackground()

	return service
}

func (s *Service) HandleMessage(m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.discord.State.User.ID {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if tokens[0] != setRemindersCommand {
		return
	}
	tokens = tokens[1:]

	if m.Message.MessageReference == nil {
		_, _ = s.discord.ChannelMessageSend(
			m.ChannelID,
			"Ответь на сообщение которое я должен тебе припомнить.",
		)
		return
	}

	if len(tokens) != 2 {
		_, _ = s.discord.ChannelMessageSend(
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
		_, _ = s.discord.ChannelMessageSend(m.ChannelID, "Бро, введи нормальное число первым аргументом, а?")
		return
	}
	duration, ok := durations[durationToken]
	if !ok {
		_, _ = s.discord.ChannelMessageSend(
			m.ChannelID,
			"Второй аргумент выбери из: seconds, minutes, hours, days, weeks, months, years.",
		)
		return
	}

	_, _ = s.discord.ChannelMessageSend(m.ChannelID, "Записала, обязательно напомню.")

	reminder := sqlite.Reminder{
		MessageID: m.MessageReference.MessageID,
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		RemindAt:  time.Now().Add(time.Duration(count * int64(duration))),
	}

	err = s.queries.AddReminder(context.Background(), sqlite.AddReminderParams{
		MessageID: reminder.MessageID,
		UserID:    reminder.UserID,
		ChannelID: reminder.ChannelID,
		RemindAt:  reminder.RemindAt,
	})
	if err != nil {
		_, _ = s.discord.ChannelMessageSend(
			m.ChannelID,
			"Ошибка при создании напоминания... :(",
		)
		return
	}

	s.reminders = append(s.reminders, reminder)
}

func (s *Service) checkInBackground() {
	ticker := time.NewTicker(checkInterval)
	go func() {
		for range ticker.C {
			s.checkReminders()
		}
	}()
}

func (s *Service) loadReminders() error {
	reminders, err := s.queries.GetReminders(context.Background())
	if err != nil {
		return err
	}

	s.reminders = reminders
	return nil
}

func (s *Service) checkReminders() {
	for i, reminder := range s.reminders {
		if reminder.RemindAt.Before(time.Now()) {
			s.remind(reminder)
			s.reminders = append(s.reminders[:i], s.reminders[i+1:]...)
			s.queries.DeleteReminder(context.Background(), reminder.ID)
		}
	}
}

func (s *Service) remind(reminder sqlite.Reminder) {
	_, _ = s.discord.ChannelMessageSendComplex(
		reminder.ChannelID,
		&discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s>, напоминаю. 🙂", reminder.UserID),
			Reference: &discordgo.MessageReference{
				MessageID: reminder.MessageID,
				ChannelID: reminder.ChannelID,
				GuildID:   s.config.GuildID,
			},
		},
	)
}
