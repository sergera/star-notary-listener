package main

import (
	"github.com/sergera/star-notary-listener/internal/event"
	"github.com/sergera/star-notary-listener/internal/logger"
)

func main() {
	logger.Setup()
	defer logger.Sync()
	event.Listen()
}
