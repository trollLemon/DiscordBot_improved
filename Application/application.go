package application

import (
	factories "bot/Core/Factories"
	"bot/Core/Interfaces"
	audio "bot/Core/Services/Audio"
	database "bot/Core/Services/Database"
	imagemanip "bot/Core/Services/ImageManip"
)

type Application struct {
	ImageApi       imagemanip.AbstractImageAPI      // Wrapper around image api
	AudioPlayer    audio.AbstractAudioPlayer        // Audio Player for voice channels
	WordDatabase   database.AbstractDatabaseService // Database Connection for random word list
	Search         Interfaces.Search                // Search API for videos
	ServiceFactory factories.ServiceFactory         // abstract factory for creating dependencies on the fly
	GuildID        string                           // Discord Guild the bot is in
}
