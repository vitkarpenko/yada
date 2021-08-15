package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Image struct {
	ID   string
	Body []byte
}

const imagesPerReactionLimit = 5

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

func (y *Yada) loadImagesInBackground() {
	y.processImages()
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for range ticker.C {
			y.processImages()
		}
	}()
}

func (y *Yada) processImages() {
	var currentLastID string

	for {
		messages, err := y.Discord.ChannelMessages(
			y.Config.ImagesChannelID,
			loadMessagesLimit,
			"",
			currentLastID,
			"",
		)
		if err != nil {
			log.Fatalln("Could not load images from image channel!", err)
		}

		y.downloadImages(messages)

		if len(messages) < loadMessagesLimit {
			break
		}

		currentLastID = messages[len(messages)-1].ID
	}
}

func (y *Yada) downloadImages(messages []*discordgo.Message) {
	for _, message := range messages {
		attachments := message.Attachments
		if len(attachments) == 0 {
			continue
		}
		images := make([]Image, len(attachments))
		for _, a := range attachments {
			images = append(images, readImageFromAttach(a))
		}
		triggerWords := strings.Split(message.Content, " ")
		y.setImagesTokens(triggerWords, images)
	}
}

func (y *Yada) setImagesTokens(triggerWords []string, images []Image) {
	for _, w := range triggerWords {
		for _, i := range images {
			y.Images[w] = i
		}
	}
}

func readImageFromAttach(a *discordgo.MessageAttachment) Image {
	response, err := http.Get(a.URL)
	if err != nil {
		return Image{}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	imageBody, _ := io.ReadAll(response.Body)
	return Image{
		ID:   a.ID,
		Body: imageBody,
	}
}
