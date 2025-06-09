package main

import (
	"bot/Application"
	"bot/Core/Commands"
	factories "bot/Core/Factories"
	"bot/Core/Wrappers"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

type Options struct {
	RegisterCommands bool
}

func loadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

}

func registerCommands(session *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands.SlashCommands))
	gid := os.Getenv("GUILD_ID")
	for i, v := range Commands.SlashCommands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, gid, v)
		if err != nil {
			log.Panic().Msgf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Registered Command %v", v.Name)
	}
}

func addCommandHandlers(session *discordgo.Session, app *application.Application) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := Commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {

			wrappedSession := Wrappers.DiscordSessionWrapper{
				Session: s,
			}

			wrappedInteractionCreate := Wrappers.InteractionCreateWrapper{
				InteractionCreate: i,
			}
			if err := h(wrappedSession, wrappedInteractionCreate, app); err != nil {
				log.Error().Err(err).Msg("Failed to execute command")
			}

		}
	})
}

func parseCommandLineArgs() *Options {

	shouldRegisterCommands := flag.Bool("register_commands", true, "register bot commands to guild")

	flag.Parse()

	return &Options{
		*shouldRegisterCommands,
	}

}

func InitializeApplication() *application.Application {

	imageApi, err := factories.CreateImageAPIService(factories.GoManip)
	if err != nil {
		log.Printf("warning, error creating image api service: %v", err)
	}

	databaseService, err := factories.CreateDatabaseService(factories.Redis)
	if err != nil {
		log.Printf("warning, error creating database service: %v", err)
	}

	classificationService, err := factories.CreateClassificationAPIService(factories.VitClassification)
	if err != nil {
		log.Printf("warning, error creating classification api service: %v", err)
	}

	return &application.Application{
		ImageApi:          imageApi,
		WordDatabase:      databaseService,
		ClassificationApi: classificationService,
	}
}

func main() {
	loadENV()
	options := parseCommandLineArgs()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal().Msg("No token provided. Set DISCORD_TOKEN in your .env file.")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal().Msgf("error creating Discord session: %v", err)
	}

	log.Info().Msg("Created a Discord Session")

	err = session.Open()
	if err != nil {
		log.Fatal().Msgf("error opening connection: %v", err)
	}

	log.Info().Msg("Connected to Discord")

	if options.RegisterCommands {
		log.Info().Msg("Registering commands...")
		registerCommands(session)
	}

	log.Info().Msg("Initializing application")
	app := InitializeApplication()

	addCommandHandlers(session, app)

	log.Info().Msg("Bot is online.")
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info().Msg("Press Ctrl+C to stop the bot")
	<-stop

}
