package Common

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Reply(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}

func ReplyImageClassification(image []byte, err error, classification string, s *discordgo.Session, i *discordgo.InteractionCreate) {

	var responseEdit *discordgo.WebhookEdit

	classificationMsg := "This is: " + classification

	if err != nil {

		errResponse := err.Error()
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
			Content: &classificationMsg,
		}
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func ReplyImage(image []byte, err error, s *discordgo.Session, i *discordgo.InteractionCreate) {

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

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func DeferReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	if err != nil {
		log.Error().Err(err).Msg("Interaction defer Response")

	}
}
