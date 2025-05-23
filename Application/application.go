package application

import (
	database "bot/Core/Services/Database"
	imagemanip "bot/Core/Services/ImageManip"
)

type Application struct {
	ImageApi     imagemanip.AbstractImageAPI      // Wrapper around image api
	WordDatabase database.AbstractDatabaseService // Database Connection for random word list
	GuildID      string                           // Discord Guild the bot is in
}
