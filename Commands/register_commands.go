package Commands

import (
	"github.com/bwmarrin/discordgo"
)

var (

	SlashCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Enqueue a youtube video and play it",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The search query",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop",
			Description: "Enqueue a youtube video and play it",
		},
		{
			Name:        "skip",
			Description: "Skip what is playing, and start the next audio in the queue",
		},
		{
			Name:        "shuffle",
			Description: "Shuffle the queue",
		},
		{
			Name:        "randomplay",
			Description: "Play a random youtube video, searched with random terms",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "Number of terms to use",
					Required:    true,
				},
			},
		},
		{
			Name:        "add",
			Description: "add text for the random play section",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The term",
					Required:    true,
				},
			},
		},
		{
			Name:        "remove",
			Description: "Remove text from the random play section",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The term",
					Required:    true,
				},
			},
		},
		{
			Name:        "show",
			Description: "Show database of random search terms",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Play(s, i)

		},
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Stop(s, i)
		},
		"shuffle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Shuffle(s, i)
		},
		"skip": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Skip(s, i)
		},

		"randomplay": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			RandomPlay(s, i)
		},
		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Add(s, i)
		},

		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Remove(s, i)
		},
		"show": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Show(s, i)
		},
	}
)
