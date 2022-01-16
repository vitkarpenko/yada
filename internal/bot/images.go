package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Body []byte

type Images struct {
	MessageID string
	Bodies    []Body
}

const imagesPerReactionLimit = 5

func (y *Yada) ReactWithImageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	words := tokenize(m.Content)
	files := y.getFilesToSend(words)

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

func (y *Yada) getFilesToSend(words []string) []*discordgo.File {
	var (
		images []Images
		files  []*discordgo.File
	)
	seenWords := make(map[string]bool)
	seenImages := make(map[string]bool)
	for _, word := range words {
		if seenWords[word] {
			continue
		}
		if image, ok := y.Images[word]; ok {
			if seenImages[image.MessageID] {
				continue
			}
			images = append(images, image)
			seenImages[image.MessageID] = true
		}
		seenWords[word] = true
	}
	for _, image := range images {
		imageToShowIndex := rand.Intn(len(image.Bodies))
		files = append(files, &discordgo.File{
			Name:        fmt.Sprintf("image_%s.gif", image.MessageID),
			ContentType: "image/gif",
			Reader:      bytes.NewReader(image.Bodies[imageToShowIndex]),
		})
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
		y.Images = make(map[string]Images)

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
		if len(attachments) == 0 || len(message.Content) == 0 {
			continue
		}

		bodies := make([]Body, len(attachments))
		for i, a := range attachments {
			bodies[i] = readImageBodyFromAttach(a)
		}

		triggerWords := strings.Split(message.Content, " ")
		images := Images{
			MessageID: message.ID,
			Bodies:    bodies,
		}
		y.setImagesTokens(triggerWords, images)
	}
}

func (y *Yada) setImagesTokens(triggerWords []string, images Images) {
	for _, w := range triggerWords {
		y.Images[w] = images
	}
}

func readImageBodyFromAttach(a *discordgo.MessageAttachment) []byte {
	response, err := http.Get(a.URL)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	imageBody, _ := io.ReadAll(response.Body)

	return imageBody
}
