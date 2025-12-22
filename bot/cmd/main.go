package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/commands"
	"github.com/trollLemon/DiscordBot/internal/common"
)

type Options struct {
	RegisterCommands bool
	PrettyPrint      bool
}



func parseCommandLineArgs() *Options {

	shouldRegisterCommands := flag.Bool("register-commands", true, "register bot commands to guild")
	prettyPrint := flag.Bool("pretty-print", false, "enable pretty printing formatting for the logs")
	flag.Parse()

	return &Options{
		*shouldRegisterCommands,
		*prettyPrint,
	}
}

func main() {
	options := parseCommandLineArgs()

	if options.PrettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Reading config from ENV")

	conf := common.GetBotConfig()

	session, err := discordgo.New("Bot " + conf.BotToken)
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
		go func() {
		    log.Info().Msg("Registering commands...")
		    commands.RegisterCommands(session)
		}()
	}

	log.Info().Msg("Initializing application")
	app := application.InitializeApplication(conf, context.Background())

	commands.AddCommandHandlers(session, app)

	log.Info().Msg("Bot is online.")
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info().Msg("Press Ctrl+C to stop the bot")
	<-stop

}
