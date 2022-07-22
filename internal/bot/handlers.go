package bot

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/vitkarpenko/yada/internal/services/images"
	"github.com/vitkarpenko/yada/internal/tokens"
	"github.com/vitkarpenko/yada/internal/utils"
)

const (
	randomImageChance = 0.01
	randomEmojiChance = 0.02
)

var (
	MehMuseEmojis     = []string{"(＃＞＜)", "<(￣ ﹌ ￣)>", "(＞﹏＜)", "ヾ( ￣O￣)ツ"}
	FineMuseEmojis    = []string{"(￣_￣)・・・", "(•ิ_•ิ)?", "┐(￣ヘ￣;)┌", "ლ(ಠ_ಠ ლ)"}
	AwesomeMuseEmojis = []string{"^ - ^", "(♡-_-♡)", "ヽ(♡‿♡)ノ", "♡( ◡‿◡ )"}
	WaifuMuseEmojis   = []string{"o(≧▽≦)o", "(❤ω❤)", "♡＼(￣▽￣)／♡", "(*♡∀♡)"}
)

func (y *Yada) AllMessagesHandler(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if ds.State.User != nil && m.Author.ID == ds.State.User.ID {
		return
	}

	y.handleImages(m)
	y.handleRandomEmoji(m)
	y.handleMuses(m)
}

func (y *Yada) handleRandomEmoji(m *discordgo.MessageCreate) {
	if utils.CheckChance(randomEmojiChance) {
		_, _ = y.Discord.ChannelMessageSendComplex(
			m.ChannelID,
			&discordgo.MessageSend{
				Content: y.Emojis.Random(),
				Reference: &discordgo.MessageReference{
					MessageID: m.Message.ID,
					ChannelID: m.ChannelID,
					GuildID:   y.Config.GuildID,
				},
			},
		)
	}
}

func (y *Yada) handleImages(m *discordgo.MessageCreate) {
	words := tokens.Tokenize(m.Content)

	var files []*discordgo.File
	if utils.CheckChance(randomImageChance) {
		files = []*discordgo.File{
			images.DiscordFileFromImage(y.Images.Random(), uuid.New().String()),
		}
	} else {
		files = y.Images.GetFilesToSend(words)
	}

	if len(files) != 0 {
		_, err := y.Discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: files,
		})
		if err != nil {
			log.Println("Couldn't send an image.", err)
		}
	}
}

func (y *Yada) handleMuses(m *discordgo.MessageCreate) {
	if m.ChannelID != y.Config.MusesChannelID || len(m.Attachments) == 0 {
		return
	}

	rating := rand.Intn(13)

	var (
		punctuation string
		emojis      []string
	)
	switch {
	case 0 <= rating && rating <= 3:
		emojis = MehMuseEmojis
		punctuation = "..."
	case 4 <= rating && rating <= 7:
		emojis = FineMuseEmojis
		punctuation = "."
	case 8 <= rating && rating <= 10:
		emojis = AwesomeMuseEmojis
		punctuation = "!"
	case 11 <= rating && rating <= 12:
		emojis = WaifuMuseEmojis
		punctuation = "!!!"
	}

	emoji := emojis[rand.Intn(len(emojis))]
	message := fmt.Sprintf("%d/12%s %s", rating, punctuation, emoji)

	_, _ = y.Discord.ChannelMessageSendComplex(
		m.ChannelID,
		&discordgo.MessageSend{
			Content: message,
			Reference: &discordgo.MessageReference{
				MessageID: m.Message.ID,
				ChannelID: m.ChannelID,
				GuildID:   y.Config.GuildID,
			},
		},
	)
}
