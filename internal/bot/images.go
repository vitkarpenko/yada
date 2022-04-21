package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type Body []byte

type Images struct {
	MessageID string
	Bodies    []Body
}

type ImageDownloadResult struct {
	images       Images
	triggerWords []string
}

const (
	imagesPerReactionLimit = 5

	randomImageChance = 0.005
	wrongImageChance  = 0.1

	downloadWorkersCount = 50
)

func (y *Yada) ReactWithImageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	words := tokenize(m.Content)

	var files []*discordgo.File
	if checkChance(randomImageChance) {
		files = []*discordgo.File{
			discordFileFromImage(y.randomImage(), uuid.New().String()),
		}
	} else {
		files = y.getFilesToSend(words)
	}

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
		if checkChance(wrongImageChance) {
			files = append(files,
				discordFileFromImage(y.randomImage(), uuid.New().String()),
			)
		} else {
			imageToShowIndex := rand.Intn(len(image.Bodies))
			imageToShow := image.Bodies[imageToShowIndex]
			files = append(files, discordFileFromImage(imageToShow, image.MessageID))
		}
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
	var wg sync.WaitGroup

	jobs := make(chan discordgo.Message, len(messages))
	for _, m := range messages {
		jobs <- *m
	}
	close(jobs)

	results := make(chan ImageDownloadResult, len(messages))
	for w := 1; w <= downloadWorkersCount; w++ {
		wg.Add(1)
		go y.downloader(jobs, results, &wg)
	}

	wg.Wait()
	close(results)

	for r := range results {
		y.setImagesTokens(r.triggerWords, r.images)
	}
}

func (y *Yada) downloader(
	jobs chan discordgo.Message,
	results chan ImageDownloadResult,
	wg *sync.WaitGroup,
) {
	for j := range jobs {
		attachments := j.Attachments
		if len(attachments) == 0 || len(j.Content) == 0 {
			return
		}

		bodies := make([]Body, len(attachments))
		for i, a := range attachments {
			bodies[i] = readImageBodyFromAttach(a)
		}

		results <- ImageDownloadResult{
			images: Images{
				MessageID: j.ID,
				Bodies:    bodies,
			},
			triggerWords: strings.Split(j.Content, " "),
		}
	}
	wg.Done()
}

func (y *Yada) setImagesTokens(triggerWords []string, images Images) {
	for _, w := range triggerWords {
		mergedBodies := append(y.Images[strings.ToLower(w)].Bodies, images.Bodies...)
		y.setBodies(w, mergedBodies)
	}
}

func (y *Yada) setBodies(w string, mergedBodies []Body) {
	imagesEntry := y.Images[strings.ToLower(w)]
	imagesEntry.Bodies = mergedBodies
	y.Images[strings.ToLower(w)] = imagesEntry
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

func (y *Yada) randomImage() Body {
	bodies := make([]Body, 0)
	for _, image := range y.Images {
		bodies = append(bodies, image.Bodies...)
	}

	return bodies[rand.Intn(len(bodies))]
}

func discordFileFromImage(image Body, imageID string) *discordgo.File {
	return &discordgo.File{
		Name:        fmt.Sprintf("image_%s.gif", imageID),
		ContentType: "image/gif",
		Reader:      bytes.NewReader(image),
	}
}
