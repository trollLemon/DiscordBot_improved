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

func ReplyImageClassification(image []byte, classification string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	classificationMsg := "This is: " + classification
	embed := &discordgo.MessageEmbed{
		Title: classificationMsg,
		Color: 0x00FF00,
	}

	responseEdit := &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:   "processed_image.png",
				Reader: bytes.NewReader(image),
			},
		},
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func ReplyGomanip(image []byte, s *discordgo.Session, i *discordgo.InteractionCreate) {
	responseEdit := &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:   "processed_image.png",
				Reader: bytes.NewReader(image),
			},
		},
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func GomanipError(s *discordgo.Session, i *discordgo.InteractionCreate, errTitle, errString string) {
	errEmbed := &discordgo.MessageEmbed{
		Title:       errTitle,
		Description: errString,
		Color:       0xFF0000,
	}

	responseEdit := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{errEmbed},
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
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

func ClassificationError(s *discordgo.Session, i *discordgo.InteractionCreate, errTitle, errString string) {
	errEmbed := &discordgo.MessageEmbed{
		Title:       errTitle,
		Description: errString,
		Color:       0xFF0000,
	}

	responseEdit := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{errEmbed},
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}
