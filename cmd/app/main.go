package main

import (
	"github.com/sergera/star-notary-listener/internal/listener"
	"github.com/sergera/star-notary-listener/internal/logger"
)

func main() {
	logger.Setup()
	defer logger.Sync()
	listener := listener.NewListener()
	listener.Listen()
}
