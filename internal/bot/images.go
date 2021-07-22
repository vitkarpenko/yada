package bot

import (
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Image struct {
	ID   string
	Body []byte
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
