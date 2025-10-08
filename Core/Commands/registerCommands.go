package Commands

import (
	application "bot/Application"
	"github.com/bwmarrin/discordgo"
)

var (
	SlashCommands = []*discordgo.ApplicationCommand{

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
			Name:        "randomfilter",
			Description: "Apply a random filter an image, for each color channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "kernel",
					Description: "length and width of the kernel (filter).",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "lowerbound",
					Description: "lowest value for random values",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "upperbound",
					Description: "highest value for random values",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "normalize",
					Description: "normalize the filter, may soften artifacts in result image",
					Required:    true,
				},
			},
		},

		{
			Name:        "invertimage",
			Description: "invert the colors of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
			},
		},
		{
			Name:        "saturateimage",
			Description: "saturate colors of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "magnitude",
					Description: "magnitude of saturation (between 0 and 100)",
					Required:    true,
				},
			},
		},
		{
			Name:        "edgedetect",
			Description: "Detect Edges in an Image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "lowerbound",
					Description: "lower bound for edge values (a good default is 100)",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "upperbound",
					Description: "upper bound for edge values (a good default is 200)",
					Required:    true,
				},
			},
		},
		{
			Name:        "dilateimage",
			Description: "enlarges objects",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "boxsize",
					Description: "how much to enlarge stuff",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "iterations",
					Description: "how many dilations to apply",
					Required:    true,
				},
			},
		},
		{
			Name:        "erodeimage",
			Description: "shrinks objects",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "boxsize",
					Description: "how much to enlarge stuff",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "iterations",
					Description: "how many dilations to apply",
					Required:    true,
				},
			},
		},
		{
			Name:        "addtext",
			Description: "add text to an image. ",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "text to paste onto image",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "fontscale",
					Description: "how big the text shall be",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "x",
					Description: "percentage of the image width between 0 and 100 (50 will be in the middle of the image",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "y",
					Description: "percentage of the image height between 0 and 100 (50 will be in the middle of the image)",
					Required:    true,
				},
			},
		},
		{
			Name:        "randomtext",
			Description: "add random text to an image (uses text database). ",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "terms",
					Description: "number of terms to get from database",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "fontscale",
					Description: "how big the text shall be",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "x",
					Description: "percentage of the image width between 0 and 100 (50 will be in the middle of the image",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "y",
					Description: "percentage of the image height between 0 and 100 (50 will be in the middle of the image)",
					Required:    true,
				},
			},
		},
		{
			Name:        "reduceimage",
			Description: "lower the quality of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "quality",
					Description: " quality (between 0 and 100) out of 100 (i.e 50/100 -> 50%)",
					Required:    true,
				},
			},
		},
		{
			Name:        "shuffleimage",
			Description: "shuffle partitions of an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to operate on",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "partitions",
					Description: "How many times to split the image up and shuffle",
					Required:    true,
				},
			},
		},
		{
			Name:        "classify",
			Description: "classify an image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "the image to classify",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error{

		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Add(s, i, a)
		},

		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Remove(s, i, a)
		},
		"show": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Show(s, i, a)
		},
		"randomfilter": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return RandomImageFilter(s, i, a)
		},
		"invertimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return InvertImage(s, i, a)
		},
		"saturateimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return SaturateImage(s, i, a)
		},
		"edgedetect": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return EdgeDetection(s, i, a)
		},
		"dilateimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Dilate(s, i, a)
		},
		"erodeimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Erode(s, i, a)
		},
		"addtext": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return AddText(s, i, a)
		},
		"randomtext": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return RandomText(s, i, a)
		},
		"reduceimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return ReduceImage(s, i, a)
		},
		"shuffleimage": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return ShuffleImage(s, i, a)
		},
		"classify": func(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
			return Classify(s, i, a)
		},
	}
)
