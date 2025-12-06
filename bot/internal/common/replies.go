package Common

import (
	"bytes"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
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

func GomanipError(s *discordgo.Session, i *discordgo.InteractionCreate, errTitle string, err error) {

	errString := "Something went wrong. Please try the command again."
	println(err.Error())
	if errors.Is(err, gomanip.ErrBadInput) {

		var userErr *gomanip.UserError

		if errors.As(err, &userErr) {
			errString = userErr.Message()
		} else {
			// fallback to error chain error string
			errString = err.Error()
		}
	}

	if errors.Is(err, gomanip.ErrRetry) {
		errString = "The request took too long; please try again."
	}

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
