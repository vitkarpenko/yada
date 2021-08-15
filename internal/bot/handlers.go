package bot

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"yada/internal/storage/postgres"

	"github.com/bwmarrin/discordgo"
)

const imagesPerReactionLimit = 5

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

func (y *Yada) ChoiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.ApplicationCommandData().Options[0].StringValue()
	words := strings.Split(message, ",")
	randIndex := rand.Intn(len(words))

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Выбирал из списка `%s` и выбрал: **%s**.",
				message,
				strings.TrimSpace(words[randIndex]),
			),
		},
	})
}

func (y *Yada) ReactWithImageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	words := tokenize(m.Content)
	files := getFilesToSend(words, y.Images)

	if len(files) != 0 {
		files = limitFilesCount(files)
		_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: files,
		})
		if err != nil {
			log.Println("Couldn't send an image.", err)
		}
	}
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

func getFilesToSend(words []string, images map[string]Image) []*discordgo.File {
	var files []*discordgo.File
	seenWords := make(map[string]bool)
	seenImages := make(map[string]bool)
	for _, word := range words {
		if seenWords[word] {
			continue
		}
		if image, ok := images[word]; ok {
			if seenImages[image.ID] {
				continue
			}
			files = append(files, &discordgo.File{
				Name:        fmt.Sprintf("image_%s.gif", image.ID),
				ContentType: "image/gif",
				Reader:      bytes.NewReader(image.Body),
			})
			seenImages[image.ID] = true
		}
		seenWords[word] = true
	}
	return files
}

func limitFilesCount(files []*discordgo.File) []*discordgo.File {
	if len(files) >= imagesPerReactionLimit {
		files = files[:imagesPerReactionLimit]
	}
	return files
}
