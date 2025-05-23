package Common

import (
	"bot/Core/Interfaces"
	"bytes"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Reply(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, text string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	}

	if err := s.InteractionRespond(i.GetInteraction(), response); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}

func ReplyImage(image []byte, err error, s Interfaces.DiscordSession, i Interfaces.DiscordInteraction) {

	var responseEdit *discordgo.WebhookEdit

	if err != nil {

		errResponse := "An error occurred: " + err.Error()
		responseEdit = &discordgo.WebhookEdit{
			Content: &errResponse,
		}
	} else {
		responseEdit = &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:   "processed_image.png",
					Reader: bytes.NewReader(image),
				},
			},
		}
	}

	if _, err := s.InteractionResponseEdit(i.GetInteraction(), responseEdit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func DeferReply(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction) {
	err := s.InteractionRespond(i.GetInteraction(), &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	if err != nil {
		log.Error().Err(err).Msg("Interaction defer Response")

	}
}
