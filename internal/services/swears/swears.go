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

	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/spelling"
)

var (
	punishmentPhrasesPath = "data/swearsPunishPhrases.txt"
	punishmentImagesPath  = "data/swearsPunishImages"
)

type Service struct {
	swears            map[string]struct{}
	punishmentImages  images.Images
	punishmentPhrases []string
}

func New() *Service {
	service := &Service{}
	service.loadSwears()
	service.loadPunishmentImages()
	service.loadPunishmentPhrases()
	return service
}

func (s *Service) IsSwear(word string) bool {
	if _, ok := s.swears[word]; ok {
		return true
	}
	return false
}

func (s *Service) PunishmentImage() images.Body {
	return s.punishmentImages.Bodies[rand.Intn(len(s.punishmentImages.Bodies))]
}

func (s *Service) PunishmentPhrase() string {
	return s.punishmentPhrases[rand.Intn(len(s.punishmentPhrases))]
}

func (s *Service) loadSwears() {
	s.swears = make(map[string]struct{})

	swearFile, err := os.Open("data/swear.tar.gz")
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
		s.addSwear(word)
	}

	fmt.Printf("Loaded up swears. %d words in dictionary!\n", len(s.swears))
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

func (s *Service) addSwear(word string) {
	s.swears[word] = struct{}{}
	edits := spelling.SimpleEdits(word, false)
	for _, e := range edits {
		s.swears[e] = struct{}{}
	}
}
