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

		{
			Name:	     "randomfilter",
			Description: "Apply a random filter an image, for each color channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},

				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "kernel",
					Description: "length and width of the kernel (filter).",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "lowerbound",
					Description: "lowest value for random values",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "upperbound",
					Description: "highest value for random values",
					Required: true,

				},			

			},

		},

		{
			Name:	     "invertimage",
			Description: "invert the colors of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},

			},
		},
		{
			Name:	     "saturateimage",
			Description: "saturate colors of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "magnitude",
					Description: "magnitude of saturation (between 0 and 100)",
					Required: true,
				},

			},
		},
		{
			Name:	     "edgedetect",
			Description: "Detect Edges in an Image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "lowerbound",
					Description: "lower bound for edge values (a good default is 100)",
					Required: true,
				},
				
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "upperbound",
					Description: "upper bound for edge values (a good default is 200)",
					Required: true,
				},
			},
		},
		{
			Name:	     "dilateimage",
			Description: "enlarges objects",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "boxsize",
					Description: "how much to enlarge stuff",
					Required: true,
				},
				
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "iterations",
					Description: "how many dilations to apply",
					Required: true,
				},
			},
		},
		{
			Name:	     "erodeimage",
			Description: "shrinks objects",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "boxsize",
					Description: "how much to enlarge stuff",
					Required: true,
				},
				
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "iterations",
					Description: "how many dilations to apply",
					Required: true,
				},
			},
		},
		{
			Name:	     "addtext",
			Description: "add text to an image. ",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "text",
					Description: "text to paste onto image",
					Required: true,
				},
				
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "fontscale",
					Description: "how big the text shall be",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "x",
					Description: "percentage of the image width between 0 and 100 (50 will be in the middle of the image",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "y",
					Description: "percentage of the image height between 0 and 100 (50 will be in the middle of the image)",
					Required: true,
				},
			},
		},

		{
			Name:	     "reduceimage",
			Description: "lower the quality of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionAttachment,
					Name: "image",
					Description: "the image to operate on",
					Required: true,

				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "quality",
					Description: " quality (between 0 and 100) out of 100 (i.e 50/100 -> 50%)",
					Required: true,
				},

			},
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
		"randomfilter": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			RandomImageFilter(s, i)
		},
		"invertimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			InvertImage(s, i)
		},
		"saturateimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			SaturateImage(s, i)
		},
		"edgedetect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			EdgeDetection(s, i)
		},
		"dilateimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Dilate(s, i)
		},
		"erodeimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Erode(s, i)
		},
		"addtext": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			AddText(s, i)
		},
		"reduceimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ReduceImage(s, i)
		},



	}
)
