package util

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
)

func GetBasicReply(content string) *discordgo.InteractionResponse {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	return response
}

func GetAttatchmentReply(image []byte) *discordgo.InteractionResponse {

	attachment := discordgo.File{
		Name:   "image.png",
		Reader: bytes.NewBuffer(image),
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{&attachment},
		},
	}

	return response

}
