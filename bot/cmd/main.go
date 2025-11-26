package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	application "github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/classification"
	"github.com/trollLemon/DiscordBot/internal/commands"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/randomwords"
)

type Options struct {
	RegisterCommands bool
	PrettyPrint      bool
}

type Config struct {
	BotToken  string
	RedisPass string

	GomanipURL        string
	ClassificationURL string
	DbDSN             string

	GomanipTimeout        time.Duration
	ClassificationTimeout time.Duration

	RandomWordsSetName string
	RedisSetNumber     int
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
			if err := h(s, i, app); err != nil {
				log.Error().Err(err).Msg("Failed to execute command")
			}

		}
	})
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

func InitializeApplication(conf *Config, ctx context.Context) *application.Application {

	gomanip := gomanip.NewGoManip(conf.GomanipURL, conf.GomanipTimeout)

	classifier := Classification.NewImageClassification(conf.ClassificationTimeout, conf.ClassificationURL, Classification.SendImageEndpoint, Classification.GetClassificationEndpoint)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.DbDSN,
		Password: conf.RedisPass,
		DB:       conf.RedisSetNumber,
	})

	redisInstance := store.NewRedisClient(ctx, redisClient, conf.RandomWordsSetName)

	randomWords := store.NewRandomWords(redisInstance)

	return &application.Application{
		Gomanip:        gomanip,
		Classification: classifier,
		RandomWords:    randomWords,
	}
}

func configure() *Config {
	conf := Config{
		BotToken:           os.Getenv("DISCORD_TOKEN"),
		RedisPass:          os.Getenv("REDIS_PASS"),
		GomanipURL:         os.Getenv("GOMANIP_URL"),
		ClassificationURL:  os.Getenv("CLASSIFICATION_URL"),
		DbDSN:              os.Getenv("DATABASE_DSN"),
		RandomWordsSetName: os.Getenv("RANDOM_WORD_SET_NAME"),
	}

	gomanipTimeout, err := time.ParseDuration(os.Getenv("GOMANIP_TIMEOUT"))
	if err != nil {
		log.Error().Msgf("failed to parse duration in provided variable GOMANIP_TIMEOUT=%s. Using default value of 30 seconds.", os.Getenv("GOMANIP_TIMEOUT"))
		gomanipTimeout = time.Second * 30
	}

	classificationTimeout, err := time.ParseDuration(os.Getenv("CLASSIFICATION_TIMEOUT"))
	if err != nil {
		log.Error().Msgf("failed to parse duration in provided variable CLASSIFICATION_TIMEOUT=%s. Using default value of 5 minutes.", os.Getenv("CLASSIFICATION_TIMEOUT"))
		gomanipTimeout = time.Minute * 5
	}

	conf.GomanipTimeout = gomanipTimeout
	conf.ClassificationTimeout = classificationTimeout

	setNum, err := strconv.Atoi(os.Getenv("REDIS_SET_NUM"))
	if err != nil {
		log.Error().Msgf("failed to parse string in provided variable REDIS_SET_NUM=%s. Using default value of 0.", os.Getenv("REDIS_SET_NUM"))
	}

	conf.RedisSetNumber = setNum

	return &conf
}

func main() {
	options := parseCommandLineArgs()

	if options.PrettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Reading config from ENV")

	conf := configure()

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
		log.Info().Msg("Registering commands...")
		registerCommands(session)
	}

	log.Info().Msg("Initializing application")
	app := InitializeApplication(conf, context.Background())

	addCommandHandlers(session, app)

	log.Info().Msg("Bot is online.")
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info().Msg("Press Ctrl+C to stop the bot")
	<-stop

}
