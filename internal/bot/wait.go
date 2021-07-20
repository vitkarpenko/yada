package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// waitUntilInterrupted waits  until CTRL-C or other term signal is received.
func waitUntilInterrupted() {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	waitChannel := make(chan os.Signal, 1)
	signal.Notify(waitChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-waitChannel
}
