package util

import (
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
