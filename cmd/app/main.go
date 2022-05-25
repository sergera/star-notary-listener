package main

import (
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/event"
	"github.com/sergera/star-notary-listener/internal/logger"
)

func main() {
	env.Init()
	logger.Init()
	defer logger.Sync()
	eth.Init()
	event.Listen()
}
