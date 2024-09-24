package Commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/dgvoice"
	"github.com/kkdai/youtube/v2"
	"github.com/jonas747/dca"
	"log"
	"os"

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
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			query := i.ApplicationCommandData().Options[0].StringValue()

			initialResponse := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Searching for '%s'...", query),
				},
			}

			if err := s.InteractionRespond(i.Interaction, initialResponse); err != nil {
				log.Printf("error responding to interaction: %v", err)
				return
			}


		},
	}
)
