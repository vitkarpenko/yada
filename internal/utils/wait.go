package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

// waitUntilInterrupted waits  until CTRL-C or other term signal is received.
func WaitUntilInterrupted() {
	log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	waitChannel := make(chan os.Signal, 1)
	signal.Notify(waitChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-waitChannel
}
