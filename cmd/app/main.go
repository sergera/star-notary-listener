package main

import (
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/event"
)

func main() {
	env.Init()
	eth.Init()
	event.Listen()
}
