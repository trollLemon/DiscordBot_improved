package application

import (
	"bot/Core/Services/Classification"
	database "bot/Core/Services/Database"
	imagemanip "bot/Core/Services/ImageManip"
)

type Application struct {
	ImageApi          imagemanip.AbstractImageAPI              // Wrapper around image api
	ClassificationApi Classification.AbstractClassificationAPI // Wrapper around classification api
	WordDatabase      database.AbstractDatabaseService         // Database wrapper for random word list
	GuildID           string                                   // Discord Guild the bot is in
}
