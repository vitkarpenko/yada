package gachi

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rs/zerolog/log"
)

const maxAutocompleteOpts = 25

type Service struct {
	dataPath string
	sounds   []string
}

func New(dataPath string) *Service {
	service := &Service{
		dataPath: dataPath,
	}
	service.setSounds()
	return service
}

func (s *Service) Handler(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Interaction.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		options := autocompleteOpts(interaction, s)
		if len(options) == 0 {
			return
		}
		sendAutocomplete(discord, interaction, options)
	case discordgo.InteractionApplicationCommand:
		fileName := interaction.ApplicationCommandData().Options[0].StringValue()
		file, err := os.ReadFile(filepath.Join(s.dataPath, fileName))
		if err != nil {
			sendSoundNotFound(discord, interaction, fileName)
			return
		}
		sendSound(discord, interaction, fileName, file)
	}
}

func sendSound(
	discord *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	fileName string, file []byte,
) {
	_ = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{discordWavFromBytes(fileName, file)},
		},
	})
}

func sendSoundNotFound(
	discord *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	fileName string,
) {
	_ = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Не могу найти такой файл: %s... 🤔", fileName),
		},
	})
}

func autocompleteOpts(
	interaction *discordgo.InteractionCreate, s *Service,
) []*discordgo.ApplicationCommandOptionChoice {
	query := interaction.ApplicationCommandData().Options[0].StringValue()
	options := s.complete(query)
	return options
}

func sendAutocomplete(
	discord *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	options []*discordgo.ApplicationCommandOptionChoice,
) {
	_ = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: options,
		},
	})
}

func discordWavFromBytes(fileName string, data []byte) *discordgo.File {
	file := &discordgo.File{
		Name:        fileName,
		ContentType: "audio/x-wav",
		Reader:      bytes.NewReader(data),
	}
	return file
}

func (s *Service) complete(query string) []*discordgo.ApplicationCommandOptionChoice {
	if len(query) == 0 {
		return nil
	}

	matches := fuzzy.FindFold(query, s.sounds)
	if len(matches) == 0 {
		return nil
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(matches))
	for i, m := range matches {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  m,
			Value: m + ".wav",
		}
	}

	if len(choices) > maxAutocompleteOpts {
		choices = choices[:maxAutocompleteOpts]
	}

	return choices
}

func (s *Service) setSounds() {
	filepath.WalkDir(
		s.dataPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				filename := filepath.Base(path)
				splitted := strings.Split(filename, ".")
				if len(splitted) != 2 || splitted[1] != "wav" {
					log.Warn().Msgf("Incorrect gachi sound file '%s', only .wav files are allowed", filename)
					return nil
				}
				s.sounds = append(s.sounds, splitted[0])
			}

			return nil
		},
	)
}
