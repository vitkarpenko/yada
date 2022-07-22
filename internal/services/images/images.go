package images

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
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"github.com/vitkarpenko/yada/internal/spelling"
	"github.com/vitkarpenko/yada/internal/utils"
)

const (
	downloadWorkersCount      = 50
	imagesPerReactionLimit    = 5
	loadMessagesLimit         = 100
	wrongImageChance          = 0.02
	redownloadTimeout         = 20 * time.Second
	cacheCleanPeriod          = 30 * time.Minute
	minWordLengthToSpellcheck = 4
)

type Service struct {
	images map[string]Images

	msgsIDsCache map[string]struct{}
	mu           sync.Mutex

	discord         *discordgo.Session
	imagesChannelID string
}

func New(discord *discordgo.Session, imagesChannelID string) *Service {
	service := &Service{
		discord:         discord,
		images:          make(map[string]Images),
		msgsIDsCache:    make(map[string]struct{}),
		imagesChannelID: imagesChannelID,
	}
	service.loadInBackground()
	service.cleanCachePeriodically()
	return service
}

type Body []byte

type Images struct {
	MessageID string
	Bodies    []Body
}

type ImageDownloadResult struct {
	images       Images
	triggerWords []string
}

func (s *Service) GetFilesToSend(words []string) []*discordgo.File {
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
		if image, ok := s.images[word]; ok {
			if seenImages[image.MessageID] {
				continue
			}
			images = append(images, image)
			seenImages[image.MessageID] = true
		}
		seenWords[word] = true
	}
	for _, image := range images {
		if utils.CheckChance(wrongImageChance) {
			files = append(files,
				DiscordFileFromImage(s.Random(), uuid.New().String()),
			)
		} else {
			imageToShowIndex := rand.Intn(len(image.Bodies))
			imageToShow := image.Bodies[imageToShowIndex]
			files = append(files, DiscordFileFromImage(imageToShow, image.MessageID))
		}
	}

	if len(files) != 0 {
		files = limitFilesCount(files)
	}

	return files
}

func limitFilesCount(files []*discordgo.File) []*discordgo.File {
	if len(files) >= imagesPerReactionLimit {
		files = files[:imagesPerReactionLimit]
	}
	return files
}

func (s *Service) loadInBackground() {
	s.processMessages()
	ticker := time.NewTicker(redownloadTimeout)
	go func() {
		for range ticker.C {
			s.mu.Lock()
			s.processMessages()
			s.mu.Unlock()
		}
	}()
}

func (s *Service) cleanCachePeriodically() {
	ticker := time.NewTicker(cacheCleanPeriod)
	go func() {
		for range ticker.C {
			s.mu.Lock()
			s.msgsIDsCache = make(map[string]struct{})
			s.images = make(map[string]Images)
			s.processMessages()
			s.mu.Unlock()
		}
	}()
}

func (s *Service) processMessages() {
	var currentLastID string

	for {
		messages, err := s.discord.ChannelMessages(
			s.imagesChannelID,
			loadMessagesLimit,
			currentLastID,
			"",
			"",
		)
		if err != nil {
			log.Fatalln("Could not load images from image channel!", err)
		}

		s.download(messages)

		if len(messages) < loadMessagesLimit {
			break
		}

		currentLastID = messages[len(messages)-1].ID
	}

	fmt.Printf("Downloaded images. %d words in dictionary!\n", len(s.images))
}

func (s *Service) download(messages []*discordgo.Message) {
	var wg sync.WaitGroup

	jobs := make(chan discordgo.Message, len(messages))
	for _, m := range messages {
		if _, ok := s.msgsIDsCache[m.ID]; ok {
			continue
		}
		s.msgsIDsCache[m.ID] = struct{}{}

		jobs <- *m
	}
	close(jobs)

	results := make(chan ImageDownloadResult, len(messages))
	for w := 1; w <= downloadWorkersCount; w++ {
		wg.Add(1)
		go s.downloader(jobs, results, &wg)
	}

	wg.Wait()
	close(results)

	for r := range results {
		s.setTokens(r.triggerWords, r.images)
	}
}

func (s *Service) downloader(
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

func (s *Service) setTokens(triggerWords []string, images Images) {
	for _, w := range triggerWords {
		w = strings.ToLower(w)

		var edits []string
		if utf8.RuneCountInString(w) < minWordLengthToSpellcheck {
			edits = []string{w}
		} else {
			edits = spelling.SimpleEdits(w)
		}

		for _, edit := range edits {
			mergedBodies := append(s.images[edit].Bodies, images.Bodies...)
			s.setBodies(edit, mergedBodies)
		}
	}
}

func (s *Service) setBodies(token string, mergedBodies []Body) {
	imagesEntry := s.images[token]
	imagesEntry.Bodies = mergedBodies
	s.images[token] = imagesEntry
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

func (s *Service) Random() Body {
	bodies := make([]Body, 0)
	for _, image := range s.images {
		bodies = append(bodies, image.Bodies...)
	}

	return bodies[rand.Intn(len(bodies))]
}

func DiscordFileFromImage(image Body, imageID string) *discordgo.File {
	return &discordgo.File{
		Name:        fmt.Sprintf("image_%s.gif", imageID),
		ContentType: "image/gif",
		Reader:      bytes.NewReader(image),
	}
}
