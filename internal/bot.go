package internal

import (
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
	"math/rand"
)

type Bot struct {
	Ctx *bot.Context
}

func (b *Bot) Help(*gateway.MessageCreateEvent) (string, error) {
	return b.Ctx.Help(), nil
}

// Choice randomly chooses something for you.
func (b *Bot) Choice(e *gateway.MessageCreateEvent, choices ...string) (string, error) {
	index := rand.Intn(len(choices))
	return choices[index], nil
}
