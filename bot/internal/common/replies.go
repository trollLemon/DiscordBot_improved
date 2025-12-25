package common

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/randomwords"
	"github.com/trollLemon/DiscordBot/internal/classification"
)


var defaultErrString = "Something went wrong, try the command again."


func Reply(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {

	errEmbed := &discordgo.MessageEmbed{
		Description: text,
		Color:       0x00FF00,
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
        Embeds: []*discordgo.MessageEmbed{errEmbed},
    },
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}

func ReplyError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {

	errEmbed := &discordgo.MessageEmbed{
		Title: "Command Failed",
		Description: err.Error(),
		Color:       0xFF0000,
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
        Embeds: []*discordgo.MessageEmbed{errEmbed},
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
	errString := defaultErrString
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

func ClassificationError(s *discordgo.Session, i *discordgo.InteractionCreate, errTitle string, err error) {
	errString := defaultErrString
	
	if errors.Is(err, classification.ErrBadFileType) {
		errString = "This command only works with PNG or JPEG images"
	}
	
	if errors.Is(err, classification.ErrRetry) {
		errString = "The command timed out, try using the command again."
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

func RandomWordsError(s *discordgo.Session, i *discordgo.InteractionCreate, errTitle, input string, err error) {
	errString := defaultErrString 

	if errors.Is(err, randomwords.ErrDuplicate) {
		errString = fmt.Sprintf("`%s` is already in the word list.", input)
	}
	if errors.Is(err, randomwords.ErrNotFound) {
		errString = fmt.Sprintf("`%s` is not in the word list.", input )
	}

	if errors.Is(err, randomwords.ErrEmpty) {
		errString = "The word list is empty."
	}

	errEmbed := &discordgo.MessageEmbed{
		Title:       errTitle,
		Description: errString,
		Color:       0xFF0000,
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
        Embeds: []*discordgo.MessageEmbed{errEmbed},
    },
	}

	if  err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}
