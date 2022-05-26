package swears

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/services/storage/sqlite"
	"github.com/vitkarpenko/yada/internal/spelling"
)

const (
	swearsPath                = "data/swears.gz"
	punishmentPhrasesPath     = "data/swearsPunishPhrases.txt"
	punishmentImagesPath      = "data/swearsPunishImages"
	uploadChunkSize           = 5000
	minWordLengthToSpellcheck = 4
)

type Service struct {
	db                *sqlite.DB
	punishmentImages  images.Images
	punishmentPhrases []string
}

func New(db *sqlite.DB) *Service {
	service := &Service{db: db}
	if db.ShouldFillSwears() {
		service.loadSwears()
	}
	service.loadPunishmentImages()
	service.loadPunishmentPhrases()
	return service
}

func (s *Service) HasSwear(phrase []string) string {
	swear, err := s.db.HasSwear(phrase)
	if err != nil {
		fmt.Println("Error while checking swearing:", err)
		return ""
	}

	return swear
}

func (s *Service) PunishmentImage() images.Body {
	return s.punishmentImages.Bodies[rand.Intn(len(s.punishmentImages.Bodies))]
}

func (s *Service) PunishmentPhrase() string {
	return s.punishmentPhrases[rand.Intn(len(s.punishmentPhrases))]
}

func (s *Service) loadSwears() {
	swears := make([]string, 0)

	swearFile, err := os.Open(swearsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer swearFile.Close()

	r, err := gzip.NewReader(swearFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	uncompressed := new(bytes.Buffer)
	io.Copy(uncompressed, r)

	scanner := bufio.NewScanner(uncompressed)
	for scanner.Scan() {
		word := scanner.Text()
		swears = append(swears, word)
	}

	var count int
	for i := 0; i < len(swears); i += uploadChunkSize {
		end := i + uploadChunkSize
		if end > len(swears) {
			end = len(swears)
		}

		edits := make([]string, 0)
		for _, w := range swears[i:end] {
			var wordEdits []string
			if utf8.RuneCountInString(w) < minWordLengthToSpellcheck {
				wordEdits = []string{w}
			} else {
				wordEdits = spelling.SimpleEdits(w)
			}

			edits = append(edits, wordEdits...)
			count += len(wordEdits)
			if count%100_000 == 0 {
				fmt.Printf("Uploaded %d count swears.\n", count)
			}
		}

		s.db.UploadSwears(edits)
	}

	fmt.Printf("Loaded up swears. %d words in dictionary!\n", count)
}

func (s *Service) loadPunishmentImages() {
	gifs := images.Images{}

	err := filepath.Walk(
		punishmentImagesPath,
		func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}

			gif, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			gifs.Bodies = append(gifs.Bodies, gif)

			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	s.punishmentImages = gifs

	fmt.Printf("Loaded up swear punishment %d images.\n", len(s.punishmentImages.Bodies))
}

func (s *Service) loadPunishmentPhrases() {
	phrases := make([]string, 0)

	f, err := os.Open(punishmentPhrasesPath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		phrases = append(phrases, scanner.Text())
	}

	s.punishmentPhrases = phrases

	fmt.Printf("Loaded up swear punishment %d phrases.\n", len(s.punishmentPhrases))
}
