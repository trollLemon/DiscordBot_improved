package main

import (
	"bot/Core/Commands"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

type Options struct {
	RegisterCommands bool
}

func loadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func registerCommands(session *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands.SlashCommands))
	gid := os.Getenv("GUILD_ID")
	for i, v := range Commands.SlashCommands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, gid, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Registered Command %v", v.Name)
	}
}

func addCommandHandlers(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := Commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
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

func main() {
	loadENV()
	options := parseCommandLineArgs()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("No token provided. Set DISCORD_TOKEN in your .env file.")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	log.Println("Created a Discord Session")

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	log.Println("Connected to Discord")

	log.Println("Getting Commands Ready")

	if options.RegisterCommands {
		registerCommands(session)
	}

	addCommandHandlers(session)

	log.Println("Commands ready")

	log.Println("Bot is online.")
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to stop the bot")
	<-stop

}
