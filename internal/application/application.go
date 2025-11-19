package application

import (
	"github.com/trollLemon/DiscordBot/internal/classification"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/randomwords"
)

type Application struct {
	Gomanip        *gomanip.GoManip
	Classification *Classification.ImageClassification
	RandomWords    *store.RandomWords
	GuildID        string
}
